// Package fsm implements a simple Finite State Machine
package fsm

import (
	"errors"
	"fmt"
)

type (
	// State of FSM
	State int

	// Input event
	Input int

	// Handler is a transition handler
	// when returns !nil transition failed, state unchanged
	Handler func() error

	// Transition record
	Transition struct {
		Input         Input
		Current, Next State
		Handler       Handler
	}

	// Transitions table
	Transitions []Transition

	state struct {
		next    State
		handler Handler
	}

	transition struct {
		currents map[State]*state
	}

	// FSM is a Finite State Machine
	FSM struct {
		State  State // current FSM state
		inputs map[Input]*transition
	}

	// InputError when call Do
	InputError struct {
		Input   Input
		Current State
	}

	// StateError when call New
	StateError struct {
		*InputError
		Next State
	}
)

// Error Input
func (e *InputError) Error() string {
	return fmt.Sprintf("unexpected transition, Input: %#v, Current: %#v",
		e.Input, e.Current)
}

// Error State
func (e *StateError) Error() string {
	return fmt.Sprintf("nondeterminictic transition, Input: %#v, Current: %#v, Next: %#v",
		e.Input, e.Current, e.Next)
}

// New returns new FSM instance
func New(transitions Transitions) (*FSM, error) {
	if len(transitions) > 0 {
		fsm := &FSM{
			State:  (transitions)[0].Current,
			inputs: make(map[Input]*transition),
		}
		for _, t := range transitions {
			if err := fsm.on(t.Input, t.Current, t.Next, t.Handler); err != nil {
				return nil, err
			}
		}
		return fsm, nil
	}
	return nil, errors.New("empty transitions")
}

func (fsm *FSM) on(input Input, current, next State, handler Handler) error {
	t, found := fsm.inputs[input]
	if !found {
		t = &transition{
			currents: make(map[State]*state),
		}
	}
	s, found := t.currents[current]
	if !found {
		s = &state{
			next:    next,
			handler: handler,
		}
	} else {
		return &StateError{
			InputError: &InputError{
				Input:   input,
				Current: current,
			},
			Next: next,
		}
	}
	t.currents[current] = s
	fsm.inputs[input] = t
	return nil
}

// Do pass input, call handler if present
func (fsm *FSM) Do(input Input) error {
	if t, found := fsm.inputs[input]; found {
		if s, found := t.currents[fsm.State]; found {
			if s.handler != nil {
				if err := s.handler(); err != nil {
					return err
				}
			}
			fsm.State = s.next
			return nil
		}
	}
	return &InputError{
		Input:   input,
		Current: fsm.State,
	}
}
