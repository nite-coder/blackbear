package connpool

import (
	"context"
	"errors"
	"net"
	"runtime"
	"sync"
	"time"
)

// Pool is the interface for a tcp connection pool
type Pool interface {
	// Get returns a connection from the pool, the number of retries and a boolean indicating if the connection is new
	// It takes a context to control the deadline and timeout
	// It also takes an address to specify which host to connect
	Get(ctx context.Context, address string) (net.Conn, int32, bool, error)
	// Put returns a connection to the pool
	Put(net.Conn) error
	// Status returns a map of addresses and their connection status
	Status() map[string]*PoolStatus
}

// PoolStatus is a struct that contains the current pool size and total connection number for an address
type PoolStatus struct {
	// PooledConns is the current number of connections in the pool
	PooledConns int32
	// OpenConns is the total number of connections for the address
	OpenConns int32
	// FreeCons is the total number of conntions that left could be opened.
	FreeConns int32
	// conns is the channel of connections for the address
	conns chan net.Conn
}

// pools is the struct that implements the Pool interface
type pools struct {
	// mu is a mutex to protect the pool map
	mu sync.Mutex
	// pool is a map of addresses and their connection status
	pool map[string]*PoolStatus
	// opts is a struct that contains the options for the pool
	opts options
}

// options is a struct that contains the options for the pool
type options struct {
	// poolSize is the size of each connection channel
	poolSize int32
	// idleConnTimeout is the duration to close idle connections
	idleConnTimeout time.Duration
	// connTimeout is the timeout for dialing a new connection
	connTimeout time.Duration
	// maxOpenConns is the maximum number of open connections per address
	maxOpenConns int32
	// maxRetry is the maximum number of retries for getting a connection
	maxRetry int32
	// readTimeout is the timeout for reading from a connection
	readTimeout time.Duration
	// writeTimeout is the timeout for writing to a connection
	writeTimeout time.Duration
}

// Option is the type for the functional options for the pool
type Option func(*options)

// WithPoolSize sets the poolSize option
func WithPoolSize(poolSize int32) Option {
	return func(opts *options) {
		opts.poolSize = poolSize
	}
}

// WithIdleConnTimeout sets the idleConnTimeout option
func WithIdleConnTimeout(idleConnTimeout time.Duration) Option {
	return func(opts *options) {
		opts.idleConnTimeout = idleConnTimeout
	}
}

// WithConnTimeout sets the connTimeout option
func WithConnTimeout(connTimeout time.Duration) Option {
	return func(opts *options) {
		opts.connTimeout = connTimeout
	}
}

// WithMaxOpenConns sets the maxOpenConns option
func WithMaxOpenConns(maxOpenConns int32) Option {
	return func(opts *options) {
		opts.maxOpenConns = maxOpenConns
	}
}

// WithMaxRetry sets the maxRetry option
func WithMaxRetry(maxRetry int32) Option {
	return func(opts *options) {
		opts.maxRetry = maxRetry
	}
}

// WithReadTimeout sets the readTimeout option
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(opts *options) {
		opts.readTimeout = readTimeout
	}
}

// WithWriteTimeout sets the writeTimeout option
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(opts *options) {
		opts.writeTimeout = writeTimeout
	}
}

// NewPools creates a new pool with the given options
func NewPools(opts ...Option) Pool {
	p := &pools{
		pool: make(map[string]*PoolStatus),
		opts: options{
			poolSize:        int32(runtime.NumCPU()), // use the number of CPUs as the default pool size
			idleConnTimeout: 1 * time.Hour,           // use 1 hour as the default idle connection timeout
			connTimeout:     5 * time.Second,         // use 5 seconds as the default connection timeout
			maxOpenConns:    0,                       // use 0 as the default maximum open connection number, meaning no limit
			maxRetry:        0,                       // use 0 as the default maximum retry number, meaning no retry
			readTimeout:     0,                       // use 0 as the default read timeout, meaning no timeout
			writeTimeout:    0,                       // use 0 as the default write timeout, meaning no timeout
		},
	}
	for _, opt := range opts {
		opt(&p.opts)
	}
	return p
}

// Get returns a connection from the pool, the number of retries and a boolean indicating if the connection is new
// It takes a context to control the deadline and timeout
// It also takes an address to specify which host to connect
func (p *pools) Get(ctx context.Context, address string) (net.Conn, int32, bool, error) {
	var conn net.Conn
	var err error
	var retry int32
	var isNew bool

	connTimeoutTimer := time.NewTimer(p.opts.connTimeout)
	if p.opts.connTimeout == 0 {
		connTimeoutTimer = time.NewTimer(365 * time.Hour * 24) // 1 year
	}
	defer connTimeoutTimer.Stop()

	for {
		if retry >= p.opts.maxRetry {
			return nil, retry, isNew, errors.New("connpool: maximum retry exceeded")
		}
		p.mu.Lock()
		cpool, ok := p.pool[address]
		if !ok {
			cpool = &PoolStatus{
				PooledConns: 0,
				OpenConns:   0,
				conns:       make(chan net.Conn, p.opts.poolSize),
			}
			p.pool[address] = cpool
		}
		p.mu.Unlock()

		if len(cpool.conns) == 0 && (cpool.OpenConns < p.opts.maxOpenConns || p.opts.maxOpenConns == 0) {
			// create new connection
			dialer := net.Dialer{Timeout: p.opts.connTimeout}
			conn, err = dialer.Dial("tcp", address)
			if err != nil {
				retry++
				continue
			}

			isNew = true
			if p.opts.readTimeout > 0 {
				_ = conn.SetReadDeadline(time.Now().Add(p.opts.readTimeout))
			}
			if p.opts.writeTimeout > 0 {
				_ = conn.SetWriteDeadline(time.Now().Add(p.opts.writeTimeout))
			}
			cpool.OpenConns++
			return &idleConn{Conn: conn, activedAt: time.Now()}, retry, isNew, nil
		}

		connTimeoutTimer.Reset(p.opts.connTimeout)

		select {
		case <-ctx.Done():
			return nil, retry, isNew, ctx.Err()
		case <-connTimeoutTimer.C:
			return nil, retry, isNew, errors.New("connpool: failed to get a connection from pool and connection timeout")
		case conn = <-cpool.conns:
			if time.Since(conn.(*idleConn).activedAt) > p.opts.idleConnTimeout {
				conn.Close()
				retry++
				continue
			}
			return conn, retry, isNew, nil
		}
	}
}

// Put returns a connection to the pool
func (p *pools) Put(conn net.Conn) error {
	if conn == nil {
		return nil
	}
	addr := conn.RemoteAddr().String()
	p.mu.Lock()
	defer p.mu.Unlock()
	cpool, found := p.pool[addr]
	if !found {
		return errors.New("connpool: unknown address")
	}

	select {
	case cpool.conns <- &idleConn{Conn: conn, activedAt: time.Now()}:
		return nil
	default:
		cpool.OpenConns--
		return conn.Close()
	}
}

// Status returns a map of addresses and their connection status
func (p *pools) Status() map[string]*PoolStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	res := make(map[string]*PoolStatus)
	for addr, cpool := range p.pool {
		pStatus := &PoolStatus{
			PooledConns: int32(len(cpool.conns)),
			OpenConns:   cpool.OpenConns,
			conns:       nil, // do not expose the channel to the outside
		}
		pStatus.FreeConns = p.opts.maxOpenConns - pStatus.OpenConns
		res[addr] = pStatus
	}
	return res
}

// idleConn is a wrapper for net.Conn that records the last active time
type idleConn struct {
	net.Conn
	activedAt time.Time
}