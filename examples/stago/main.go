package main

import (
	"log"

	"github.com/nenavizhuleto/sadm"
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

	locked := s.NewState(LOCKED)
	locked.Transition(PUSH, LOCKED)
	locked.Action(func(ctx *Turnstile, message Message) {
		switch message {
		case COIN:
			ctx.Count++
		case PUSH:
			if ctx.Count > 0 {
				ctx.Count--
			}
		}

	})
	locked.Transition(COIN, UNLOCKED)

	unlocked := s.NewState(UNLOCKED)
	unlocked.Transition(PUSH, LOCKED)
	unlocked.Transition(COIN, UNLOCKED)
	unlocked.Action(func(ctx *Turnstile, message Message) {
		switch message {
		case COIN:
			ctx.Count++
		case PUSH:
			if ctx.Count > 0 {
				ctx.Count--
			}
		}

	})

	unlocked.Condition(func(t Turnstile) bool {
		return t.Count > 0
	})

	unlocked.Decision(func(t Turnstile) bool {
		return t.Count > 5
	}, func(ctx *Turnstile, message Message) {
		ctx.Count = 0
	})

	unlocked.Fallback(func(ctx *Turnstile, message Message) {
	})

	adm := sadm.New("turnstile")
	adm.AddCommand("coin", "insert coin into turnstile", func(c *sadm.Connection) error {
		s.Forward(COIN)
		return c.Println(s.State(), s.Context())
	})

	adm.AddCommand("push", "push turnstile", func(c *sadm.Connection) error {
		s.Forward(PUSH)
		return c.Println(s.State(), s.Context())
	})

	adm.Listen(":3999")
}
