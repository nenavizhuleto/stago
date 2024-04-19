package main

import (
	"log"

	"github.com/nenavizhuleto/stago/automaton"
)

const (
	coin = "coin"
	push = "push"
)

type Context struct {
	Coins int
}

func main() {
	ctx := Context{Coins: 0}
	unlocked := automaton.NewStateBuilder(&ctx)
	locked := automaton.NewStateBuilder(&ctx)

	locked.OnEnter(func(ctx *Context) {
		log.Println("enter locked")
	})

	unlocked.OnEnter(func(ctx *Context) {
		log.Println("enter unlocked")
	})

	locked_state := locked.Build()
	unlocked_state := unlocked.Build()

	unlocked.OnNext(func(ctx *Context, event any) automaton.State {
		log.Println(ctx, event)
		switch event {
		case push:
			ctx.Coins = max(0, ctx.Coins-1)
			if ctx.Coins == 0 {
				return locked_state
			} else {
				return nil
			}
		case coin:
			ctx.Coins++
		}
		return nil
	})

	locked.OnNext(func(ctx *Context, event any) automaton.State {
		log.Println(ctx, event)
		switch event {
		case push:
			return nil
		case coin:
			ctx.Coins++
			return unlocked_state
		}
		return nil
	}).Build()

	a := automaton.New(locked_state)

	a.Forward(push)
	a.Forward(coin)
	a.Forward(push)
	a.Forward(coin)
	a.Forward(push)
	a.Forward(coin)
	a.Forward(coin)
	a.Forward(push)
	a.Forward(push)
	a.Forward(coin)
	a.Forward(push)
	a.Forward(push)

}
