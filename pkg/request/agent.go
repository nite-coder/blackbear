package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// httpClient should be kept for reuse purpose
	_httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			// disable "G402 (CWE-295): TLS MinVersion too low. (Confidence: HIGH, Severity: HIGH)"
			// #nosec G402
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	_timeout = 30 * time.Second

	// ErrTimeout means http request have been timeout
	ErrTimeout = errors.New("request: request timeout")
)

// Agent the main struct to handle all http requests
type Agent struct {
	client *http.Client
	err    error

	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

func newAgentWithClient(client *http.Client) Agent {
	agent := Agent{
		client:  client,
		Headers: map[string]string{},
	}
	agent.Headers["Accept"] = "application/json"
	agent.Timeout = _timeout
	return agent
}

func (a Agent) getTransport() *http.Transport {
	trans, _ := a.client.Transport.(*http.Transport)
	return trans
}

// SetMethod return Agent that uses HTTP method with target URL
func (a Agent) SetMethod(method, targetURL string) Agent {
	a.Method = strings.ToUpper(method)

	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// SetClient return Agent with target URL
func (a Agent) SetClient(client *http.Client) Agent {
	a.client = client
	return a
}

// GET return Agent that uses HTTP GET method with target URL
func (a Agent) GET(targetURL string) Agent {
	a.Method = "GET"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// POST return Agent that uses HTTP POST method with target URL
func (a Agent) POST(targetURL string) Agent {
	a.Method = "POST"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// PUT return Agent that uses HTTP PUT method with target URL
func (a Agent) PUT(targetURL string) Agent {
	a.Method = "PUT"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// DELETE return Agent that uses HTTP PUT method with target URL
func (a Agent) DELETE(targetURL string) Agent {
	a.Method = "DELETE"
	_, err := url.Parse(targetURL)
	if err != nil {
		a.err = err
	}
	a.URL = targetURL
	return a
}

// Header that set HTTP header to agent
func (a Agent) Header(key, val string) Agent {
	newHeader := map[string]string{}

	if a.Headers != nil {
		for k, val := range a.Headers {
			newHeader[k] = val
		}
	}

	newHeader[key] = val
	a.Headers = newHeader

	return a
}

// SetTimeout set timeout for agent.  The default value is 30 seconds.
func (a Agent) SetTimeout(timeout time.Duration) Agent {
	if timeout > 0 {
		a.Timeout = timeout
	}
	return a
}

// SetProxyURL set the simple proxy with fixed proxy url
func (a Agent) SetProxyURL(proxyURL string) Agent {
	trans := a.getTransport()
	if trans == nil {
		a.err = errors.New("request: no transport")
	}

	u, err := url.Parse(proxyURL)
	if err != nil {
		a.err = err
	}

	p := http.ProxyURL(u)
	if p != nil {
		trans.Proxy = p
	}

	return a
}

// SendBytes send bytes to target URL
func (a Agent) SendBytes(bytes []byte) Agent {
	a.Body = bytes
	return a
}

// SendJSON send json to target URL
func (a Agent) SendJSON(v interface{}) Agent {
	newAgent := a.Header("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		newAgent.err = err
	}
	return newAgent.SendBytes(b)
}

// SendXML send json to target URL
func (a Agent) SendXML(v interface{}) Agent {
	newAgent := a.Header("Content-Type", "application/xml")
	b, err := xml.Marshal(v)
	if err != nil {
		newAgent.err = err
	}
	return newAgent.SendBytes(b)
}

// Send send string to target URL
func (a Agent) Send(body string) Agent {
	newAgent := a.Header("Content-Type", "application/x-www-form-urlencoded")
	return newAgent.SendBytes([]byte(body))
}


// EndCtx executes the EndCtx function.
//
// It takes a context.Context parameter and returns a *Response and an error.
func (a Agent) EndCtx(ctx context.Context) (*Response, error) {
	return a.execute(ctx)
}

// End start execute agent
func (a Agent) End() (*Response, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if a.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, a.Timeout)
	}
	defer cancel()
	return a.execute(ctx)
}


// execute executes the Agent's request and returns the response.
//
// It takes a context.Context as the parameter and returns a *Response and an error.
func (a Agent) execute(ctx context.Context) (*Response, error) {
	if a.err != nil {
		return nil, a.err
	}

	// create new request
	url := a.URL
	outReq, err := http.NewRequest(a.Method, url, bytes.NewReader(a.Body))
	if err != nil {
		return nil, err
	}

	// copy Header
	for k, val := range a.Headers {
		outReq.Header.Add(k, val)
	}

	// send to target
	resp, err := a.client.Do(outReq.WithContext(ctx))
	if err != nil {
		var errTimeout net.Error
		if errors.As(err, &errTimeout) && errTimeout.Timeout() {
			return nil, ErrTimeout
		}
		return nil, err
	}
	defer func() {
		_ = respClose(resp.Body)
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	result.setResp(resp)
	result.Body = body

	if resp.StatusCode >= 200 && resp.StatusCode < 300 || resp.StatusCode == 304 {
		result.OK = true
	}

	return result, nil
}

func respClose(body io.ReadCloser) error {
	if body == nil {
		return nil
	}
	if _, err := io.Copy(io.Discard, body); err != nil {
		return err
	}
	return body.Close()
}
