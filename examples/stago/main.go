package main

import (
	"github.com/nenavizhuleto/stago"
)

type Turnstile struct {
	Locked bool
	Count  int
}

type State string

const (
	LOCKED   = State("LOCKED")
	UNLOCKED = State("UNLOCKED")
)

type Message string

const (
	PUSH = Message("PUSH")
	COIN = Message("COIN")
)

type TurnstileMessage struct {
	Push   bool
	Amount int
}

func main() {
	s := stago.New[Turnstile, State, Message](Turnstile{}, LOCKED)

	{
		state := s.NewState(LOCKED)
		state.Transition(PUSH, LOCKED)
		state.Transition(COIN, UNLOCKED)

		state.Action(func(ctx *Turnstile, message Message) {
			switch message {
			case COIN:
				ctx.Count++
			case PUSH:
				if ctx.Count > 0 {
					ctx.Count--
				}
			}
		})
	}

	{
		state := s.NewState(UNLOCKED)
		state.Transition(PUSH, LOCKED)
		state.Transition(COIN, UNLOCKED)

		state.Action(func(ctx *Turnstile, message Message) {
			switch message {
			case COIN:
				ctx.Count++
			case PUSH:
				if ctx.Count > 0 {
					ctx.Count--
				}
			}

		})

		state.Condition(func(t Turnstile) bool {
			return t.Count > 0
		})

		state.Decision(func(t Turnstile) bool {
			return t.Count > 5
		}, func(ctx *Turnstile, message Message) {
			ctx.Count = 0
		})

		state.Fallback(func(ctx *Turnstile, message Message) {})
	}
}
