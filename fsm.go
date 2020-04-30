// Package fsm implements a simple Finite State Machine
package fsm

import (
	"fmt"
)

type (
	// State of FSM
	State int

	// Input event
	Input int

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

	// Handler is a transition handler
	// when returns !nil transition failed, state unchanged
	Handler func() error

	// InputError when call Do
	InputError struct {
		Input   Input
		Current State
	}

	// StateError when call On (panic!)
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
func New(initial State) *FSM {
	return &FSM{
		State:  initial,
		inputs: make(map[Input]*transition),
	}
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
			next: next,
		}
		if handler != nil {
			s.handler = handler
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

// On defines transition, panic if nondeterminictic
func (fsm *FSM) On(input Input, current, next State, handler ...Handler) *FSM {
	var err error
	if len(handler) == 1 {
		err = fsm.on(input, current, next, handler[0])
	} else {
		err = fsm.on(input, current, next, nil)
	}
	if err != nil {
		panic(err)
	}
	return fsm
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
