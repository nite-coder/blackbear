package fsm

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrTransitionNotFound = errors.New("fsm: transition was not found")
	ErrInvalidState       = errors.New("fsm: inappropriate in current state")
)

type Event struct {
	StateMachine *StateMachine
	Transitions  string
	From         string
	To           string
	Args         []interface{}
}

type Handler func(ctx context.Context, e *Event) error

type StateMachine struct {
	mu          sync.RWMutex
	state       string
	transitions []*Transition

	beforeTransitionHandler Handler
	afterTransitionHandler  Handler
	stateHandlers           map[string]Handler
}

func New(initialState string) *StateMachine {
	return &StateMachine{
		state:         initialState,
		transitions:   []*Transition{},
		stateHandlers: map[string]Handler{},
	}
}

func (sm *StateMachine) findTransition(transition string) (*Transition, error) {
	var result *Transition
	for _, t := range sm.transitions {
		if t.name == transition {
			result = t
		}
	}

	if result == nil {
		return nil, ErrTransitionNotFound
	}

	return result, nil
}

func (sm *StateMachine) findStateHandler(key string) Handler {
	if h, ok := sm.stateHandlers[key]; ok {
		return h
	}

	return nil
}

func (sm *StateMachine) Is(sate string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.state == sate
}

func (sm *StateMachine) Can(transition string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	t, err := sm.findTransition(transition)
	if err != nil {
		return false
	}

	return contains(t.from, sm.state)
}

func (sm *StateMachine) Cannot(transition string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	t, err := sm.findTransition(transition)
	if err != nil {
		return true
	}

	return !contains(t.from, sm.state)
}

func (sm *StateMachine) State() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.state
}

func (sm *StateMachine) SetState(state string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state = state
}

func (sm *StateMachine) Trigger(ctx context.Context, transition string, args ...interface{}) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	t, err := sm.findTransition(transition)
	if err != nil {
		return err
	}

	isRightState := contains(t.from, sm.state)
	if !isRightState {
		return ErrInvalidState
	}

	// trigger
	originalState := sm.state
	evt := Event{
		StateMachine: sm,
		Transitions:  transition,
		From:         originalState,
		To:           t.to,
		Args:         args,
	}

	if sm.beforeTransitionHandler != nil {
		err = sm.beforeTransitionHandler(ctx, &evt)
		if err != nil {
			return err
		}
	}

	if t.beforeHandler != nil {
		err = t.beforeHandler(ctx, &evt)
		if err != nil {
			return err
		}
	}

	h := sm.findStateHandler("leave_" + t.to)
	if h != nil {
		err = h(ctx, &evt)
		if err != nil {
			return err
		}
	}

	// change state
	sm.state = t.to

	h = sm.findStateHandler("enter_" + t.to)
	if h != nil {
		err = h(ctx, &evt)
		if err != nil {
			return err
		}
	}

	if t.afterHandler != nil {
		err = t.afterHandler(ctx, &evt)
		if err != nil {
			return err
		}
	}

	if sm.afterTransitionHandler != nil {
		err = sm.afterTransitionHandler(ctx, &evt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *StateMachine) Transition(name string) *Transition {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	t := Transition{
		stateMachine: sm,
		name:         name,
	}

	sm.transitions = append(sm.transitions, &t)
	return &t
}

type Transition struct {
	stateMachine *StateMachine
	name         string
	from         []string
	to           string

	beforeHandler Handler
	afterHandler  Handler
}

func (t *Transition) From(states ...string) *Transition {
	t.from = states
	return t
}

func (t *Transition) To(state string) *Transition {
	t.to = state
	return t
}

func (t *Transition) Before(h Handler) *Transition {
	t.beforeHandler = h
	return t
}

func (t *Transition) After(h Handler) *Transition {
	t.afterHandler = h
	return t
}

func (t *Transition) Transition(name string) *Transition {
	return t.stateMachine.Transition(name)
}

func (t *Transition) StateMachine() *StateMachine {
	return t.stateMachine
}

func (t *Transition) BeforeTransition(h Handler) *Transition {
	t.stateMachine.beforeTransitionHandler = h
	return t
}

func (t *Transition) AfterTransition(h Handler) *Transition {
	t.stateMachine.afterTransitionHandler = h
	return t
}

func (t *Transition) LeaveState(state string, h Handler) *Transition {
	t.stateMachine.mu.Lock()
	defer t.stateMachine.mu.Unlock()

	key := "leave_" + state
	t.stateMachine.stateHandlers[key] = h
	return t
}

func (t *Transition) EnterState(state string, h Handler) *Transition {
	t.stateMachine.mu.Lock()
	defer t.stateMachine.mu.Unlock()

	key := "enter_" + state
	t.stateMachine.stateHandlers[key] = h
	return t
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
