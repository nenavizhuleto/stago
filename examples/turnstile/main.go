package main

import (
	"fmt"

	"github.com/nenavizhuleto/stago/fsm"
)

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

// 	[State-transition table](https://en.wikipedia.org/wiki/Finite-state_machine)
//	|-----------------------|---------------|-----------------------|
//	|	State/Msg	|	PUSH	|	COIN		|
//	|-----------------------|---------------|-----------------------|
//	|	LOCKED		|	LOCKED	|	UNLOCKED	|
//	|	UNLOCKED	|	LOCKED	|	UNLOCKED	|
//	|-----------------------|---------------|-----------------------|

var TABLE = fsm.TT[State, Message]{
	LOCKED: {
		PUSH: LOCKED,
		COIN: UNLOCKED,
	},
	UNLOCKED: {
		PUSH: LOCKED,
		COIN: UNLOCKED,
	},
}

func main() {
	turnstile := fsm.New(LOCKED, TABLE)

	messages := []Message{
		PUSH, // LOCKED
		PUSH, // LOCKED
		COIN, // UNLOCKED
		COIN, // UNLOCKED
		PUSH, // LOCKED
	}

	for _, message := range messages {
		turnstile.Next(message)
		fmt.Println(message, " -> ", turnstile.State())
	}
}
