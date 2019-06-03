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

	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.IsType(t, &StateError{}, r)
		assert.Equal(t, push, r.(*StateError).Input)
		assert.Equal(t, locked, r.(*StateError).Current)
		assert.Equal(t, unlocked, r.(*StateError).Next)
	}()

	New(locked).
		On(coin, locked, unlocked).
		On(coin, unlocked, unlocked).
		On(push, locked, locked).   // comment to try unexpected input
		On(push, locked, unlocked). // uncomment to try nondeterministic transition
		On(push, unlocked, locked)
}

func TestInput(t *testing.T) {
	turnstile := New(locked).
		On(coin, locked, unlocked).
		On(coin, unlocked, unlocked).
		On(push, locked, locked). // comment to try unexpected input
		// On(push, locked, unlocked). // uncomment to try nondeterministic transition
		On(push, unlocked, locked)

	assert.Equal(t, locked, turnstile.State)

	err := turnstile.Do(coin)
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
	turnstile := New(locked).
		On(coin, locked, unlocked).
		On(coin, unlocked, unlocked).
		// On(push, locked, locked).  // comment to try unexpected input
		// On(push, locked, unlocked). // uncomment to try nondeterministic transition
		On(push, unlocked, locked)

	assert.Equal(t, locked, turnstile.State)

	err := turnstile.Do(push)
	assert.NotNil(t, err)
	assert.IsType(t, &InputError{}, err)
	assert.Equal(t, push, err.(*InputError).Input)
	assert.Equal(t, locked, err.(*InputError).Current)

	assert.Equal(t, locked, turnstile.State)
}

func TestTransition(t *testing.T) {
	human := 0

	turnstile := New(locked).
		On(coin, locked, unlocked).
		On(coin, unlocked, unlocked).
		On(push, locked, locked). // comment to try unexpected input
		// On(push, locked, unlocked). // uncomment to try nondeterministic transition
		On(push, unlocked, locked, func() error {
			human++
			return nil
		})

	turnstile.Do(push)
	turnstile.Do(coin)
	turnstile.Do(push)

	assert.Equal(t, 1, human)
}

func TestTransitionFailed(t *testing.T) {

	turnstile := New(locked).
		On(coin, locked, unlocked, func() error {
			return errors.New("invalid coin")
		}).
		On(coin, unlocked, unlocked).
		On(push, locked, locked). // comment to try unexpected input
		// On(push, locked, unlocked). // uncomment to try nondeterministic transition
		On(push, unlocked, locked)

	err := turnstile.Do(coin)
	assert.NotNil(t, err)
	assert.Equal(t, locked, turnstile.State)
}
