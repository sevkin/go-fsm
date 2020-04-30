[![LICENSE](https://img.shields.io/github/license/sevkin/go-fsm.svg?color=orange)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/sevkin/go-fsm)](https://goreportcard.com/report/github.com/sevkin/go-fsm)
[![Codecov](https://img.shields.io/codecov/c/gh/sevkin/go-fsm)](https://codecov.io/gh/sevkin/go-fsm)
[![Travis (.org)](https://img.shields.io/travis/sevkin/go-fsm)](https://travis-ci.org/sevkin/go-fsm)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://godoc.org/github.com/sevkin/go-fsm)

# Finite State Machine library for Golang

1st [small piece of theory](https://en.wikipedia.org/wiki/Finite-state_machine)

## Usage on the example of our favorite turnstile

![turnstile diagram](example/turnstile.svg)

```golang
const (
    locked fsm.State = iota
    unlocked

    coin fsm.Input = iota
    push
)

turnstile, err := fsm.New(fsm.Transitions{
    // 1st Current => Initial FSM State
    {Input: coin, Current: locked, Next: unlocked, Handler: func() error {
        // return errors.New("invalid coin")
        return nil
    }},
    {coin, unlocked, unlocked, nil},
    {push, locked, locked, nil},   // comment to try unexpected input
    // {push, locked, unlocked, nil}, // uncomment to try nondeterministic transition
    {push, unlocked, locked, func() error {
        human++
        return nil
    }},
})

err = turnstile.Do(coin)
turnstile.Do(push)

fmt.Println(human, turnstile.State)
```

### Another known FSM implementations

- [github.com/rynorris/fsm](https://github.com/rynorris/fsm)
- [github.com/ryanfaerman/fsm](https://github.com/ryanfaerman/fsm)
- [github.com/looplab/fsm](https://github.com/looplab/fsm)
