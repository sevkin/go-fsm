package main

import (
	"fmt"

	"github.com/sevkin/fsm"
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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	turnstile := fsm.New(locked).
		On(coin, locked, unlocked, open).
		On(coin, unlocked, unlocked).
		On(push, locked, locked). // comment to try unexpected input
		// On(push, locked, unlocked). // uncomment to try nondeterministic transition
		On(push, unlocked, locked, close)

	turnstile.Do(coin)
	turnstile.Do(coin)
	turnstile.Do(push)
	if err := turnstile.Do(push); err != nil {
		fmt.Println(err)
	}
}
