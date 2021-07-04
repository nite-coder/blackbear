package fsm

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateMachine(t *testing.T) {
	stateMachine := New("initial").
		Transition("order_created").From("initial").To("created").
		Transition("order_paid").From("created").To("paid").
		StateMachine()

	assert.Equal(t, "initial", stateMachine.State())
	assert.Equal(t, true, stateMachine.Is("initial"))
	assert.Equal(t, true, stateMachine.Can("order_created"))
	assert.Equal(t, true, stateMachine.Cannot("order_paid"))

	stateMachine.SetState("none")
	assert.Equal(t, "none", stateMachine.State())
}

func TestTransition(t *testing.T) {
	ctx := context.Background()

	stateMachine := New("initial").
		Transition("order_created").From("initial").To("created").
		Transition("order_paid").From("created").To("paid").
		StateMachine()

	t.Run("transition success", func(t *testing.T) {
		err := stateMachine.Trigger(ctx, "order_created")
		require.NoError(t, err)

		assert.Equal(t, "created", stateMachine.State())
	})

	t.Run("transition not found", func(t *testing.T) {
		err := stateMachine.Trigger(ctx, "abc")
		assert.Equal(t, true, errors.Is(err, ErrTransitionNotFound))
	})

	t.Run("current state is incorrect", func(t *testing.T) {
		stateMachine.SetState("none")
		err := stateMachine.Trigger(ctx, "order_paid")
		assert.Equal(t, true, errors.Is(err, ErrInvalidState))
	})
}

func TestLifeCycleEventOrder(t *testing.T) {
	ctx := context.Background()
	var isBeforeTransition, isAfterTransition, isLeaveState, isEnterState, isBefore, isAfter bool

	stateMachine := New("initial").
		Transition("order_created").From("initial").To("created").
		Before(func(ctx context.Context, e *Event) error {
			if isBeforeTransition {
				isBefore = true
			}
			return nil
		}).
		After(func(ctx context.Context, e *Event) error {
			if isEnterState {
				isAfter = true
			}
			return nil
		}).
		Transition("order_paid").From("created").To("paid").
		BeforeTransition(func(ctx context.Context, e *Event) error {
			isBeforeTransition = true
			return nil
		}).
		LeaveState("initial", func(ctx context.Context, e *Event) error {
			if isBefore {
				isLeaveState = true
			}
			return nil
		}).
		EnterState("created", func(ctx context.Context, e *Event) error {
			if isLeaveState {
				isEnterState = true
			}
			return nil
		}).
		AfterTransition(func(ctx context.Context, e *Event) error {
			if isAfter && isAfterTransition == false {
				isAfterTransition = true
			}
			return nil
		}).
		StateMachine()

	err := stateMachine.Trigger(ctx, "order_created")
	require.NoError(t, err)

	assert.Equal(t, "created", stateMachine.State())
}

func TestTrigger(t *testing.T) {
	ctx := context.Background()

	stateMachine := New("initial").
		Transition("order_created").From("initial").To("created").
		Before(func(ctx context.Context, e *Event) error {
			assert.Equal(t, "orderID", e.Args[0])
			return nil
		}).
		StateMachine()

	err := stateMachine.Trigger(ctx, "order_created", "orderID")
	require.NoError(t, err)
}

func TestTriggerWithError(t *testing.T) {
	ctx := context.Background()

	stateMachine := New("initial").
		Transition("order_created").From("initial").To("created").
		Before(func(ctx context.Context, e *Event) error {
			return errors.New("something bad happened")
		}).
		StateMachine()

	err := stateMachine.Trigger(ctx, "order_created")
	assert.Equal(t, "something bad happened", err.Error())

	assert.Equal(t, "initial", stateMachine.State())
}
