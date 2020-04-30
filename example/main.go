package main

import (
	"fmt"
	"log"

	"github.com/sevkin/go-fsm"
)

func open() error {
	// return errors.New("invalid coin") // state unchanged
	fmt.Println("Unlocks the turnstile so that the customer can push through.")
	return nil
}

func close() error {
	fmt.Println("When the customer has pushed through, locks the turnstile.")
	// human++
	return nil
}

func main() {
	const (
		locked fsm.State = iota
		unlocked
	)
	const (
		coin fsm.Input = iota
		push
	)

	turnstile, err := fsm.New(fsm.Transitions{
		{Input: coin, Current: locked, Next: unlocked, Handler: open},
		{Input: coin, Current: unlocked, Next: unlocked, Handler: nil},
		{Input: push, Current: locked, Next: locked, Handler: nil}, // comment to try unexpected input
		// {push, locked, unlocked, nil},                              // uncomment to try nondeterministic transition
		{Input: push, Current: unlocked, Next: locked, Handler: close},
	})

	if err != nil {
		log.Fatal(err)
	}

	turnstile.Do(coin)
	turnstile.Do(coin)
	turnstile.Do(push)
	if err := turnstile.Do(push); err != nil {
		fmt.Println(err)
	}
}
