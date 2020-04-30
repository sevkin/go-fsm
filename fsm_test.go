package fsm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locked State = iota
	unlocked
	// )
	// const (
	coin Input = iota
	push
)

func TestNonDeterministic(t *testing.T) {
	turnstile, err := New(Transitions{
		{coin, locked, unlocked, nil},
		{coin, unlocked, unlocked, nil},
		{push, locked, locked, nil},   // comment to try unexpected input
		{push, locked, unlocked, nil}, // uncomment to try nondeterministic transition
		{push, unlocked, locked, nil},
	})

	assert.Nil(t, turnstile)

	assert.NotNil(t, err)
	assert.IsType(t, &StateError{}, err)
	if serr, ok := err.(*StateError); ok {
		assert.Equal(t, push, serr.Input)
		assert.Equal(t, locked, serr.Current)
		assert.Equal(t, unlocked, serr.Next)
	}
}

func TestInput(t *testing.T) {
	turnstile, err := New(Transitions{
		{coin, locked, unlocked, nil},
		{coin, unlocked, unlocked, nil},
		{push, locked, locked, nil}, // comment to try unexpected input
		// {push, locked, unlocked, nil}, // uncomment to try nondeterministic transition
		{push, unlocked, locked, nil},
	})

	assert.Nil(t, err)

	assert.Equal(t, locked, turnstile.State)

	err = turnstile.Do(coin)
	assert.Nil(t, err)
	assert.Equal(t, unlocked, turnstile.State)

	err = turnstile.Do(coin)
	assert.Nil(t, err)
	assert.Equal(t, unlocked, turnstile.State)

	err = turnstile.Do(push)
	assert.Nil(t, err)
	assert.Equal(t, locked, turnstile.State)

	err = turnstile.Do(push)
	assert.Nil(t, err)
	assert.Equal(t, locked, turnstile.State)
}

func TestInputUnexpected(t *testing.T) {
	turnstile, err := New(Transitions{
		{coin, locked, unlocked, nil},
		{coin, unlocked, unlocked, nil},
		// {push, locked, locked, nil}, // comment to try unexpected input
		// {push, locked, unlocked, nil}, // uncomment to try nondeterministic transition
		{push, unlocked, locked, nil},
	})

	assert.Nil(t, err)

	assert.Equal(t, locked, turnstile.State)

	err = turnstile.Do(push)
	assert.NotNil(t, err)
	assert.IsType(t, &InputError{}, err)
	assert.Equal(t, push, err.(*InputError).Input)
	assert.Equal(t, locked, err.(*InputError).Current)

	assert.Equal(t, locked, turnstile.State)
}

func TestTransition(t *testing.T) {
	human := 0

	turnstile, err := New(Transitions{
		{coin, locked, unlocked, nil},
		{coin, unlocked, unlocked, nil},
		{push, locked, locked, nil}, // comment to try unexpected input
		// {push, locked, unlocked, nil}, // uncomment to try nondeterministic transition
		{push, unlocked, locked, func() error {
			human++
			return nil
		}},
	})

	assert.Nil(t, err)

	turnstile.Do(push)
	turnstile.Do(coin)
	turnstile.Do(push)

	assert.Equal(t, 1, human)
}

func TestTransitionFailed(t *testing.T) {

	turnstile, err := New(Transitions{
		{coin, locked, unlocked, func() error {
			return errors.New("invalid coin")
		}},
		{coin, unlocked, unlocked, nil},
		{push, locked, locked, nil}, // comment to try unexpected input
		// {push, locked, unlocked, nil}, // uncomment to try nondeterministic transition
		{push, unlocked, locked, nil},
	})

	assert.Nil(t, err)

	err = turnstile.Do(coin)
	assert.NotNil(t, err)
	assert.Equal(t, locked, turnstile.State)
}
