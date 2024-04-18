package stago

import (
	"github.com/nenavizhuleto/stago/fsm"
)

type Condition[C any] func(C) bool
type Action[C any, M comparable] func(ctx *C, message M)

type Decision[C any, M comparable] struct {
	Condition Condition[C]
	Action    Action[C, M]
}

type Stago[C any, S, M comparable] struct {
	context C
	fsm     *fsm.FSM[S, M]

	decisions  map[S][]Decision[C, M]
	conditions map[S]Condition[C]
	actions    map[M][]Action[C, M]
	fallbacks  map[S]Action[C, M]
}

func New[C any, S, M comparable](ctx C, state S) *Stago[C, S, M] {
	s := &Stago[C, S, M]{
		context:    ctx,
		fsm:        fsm.New(state, fsm.TT[S, M]{}),
		decisions:  make(map[S][]Decision[C, M]),
		conditions: make(map[S]Condition[C]),
		actions:    make(map[M][]Action[C, M]),
		fallbacks:  make(map[S]Action[C, M]),
	}

	return s
}

type State[C any, S, M comparable] struct {
	Transition func(message M, next S)
	Condition  func(condition Condition[C])
	Decision   func(cond Condition[C], action Action[C, M])
	Fallback   func(action Action[C, M])
}

func (d *Stago[C, S, M]) AddDecision(state S, decision Decision[C, M]) {
	decisions, _ := d.decisions[state]
	d.decisions[state] = append(decisions, decision)
}

func (d *Stago[C, S, M]) AddFallback(state S, fallback Action[C, M]) {
	d.fallbacks[state] = fallback
}

func (d *Stago[C, S, M]) AddCondition(state S, condition Condition[C]) {
	d.conditions[state] = condition
}

func (d *Stago[C, S, M]) AddAction(message M, action Action[C, M]) {
	d.actions[message] = append(d.actions[message], action)
}

func (d *Stago[C, S, M]) NewState(state S) State[C, S, M] {
	return State[C, S, M]{
		Transition: func(message M, next S) {
			d.fsm.Table().AddTransition(state, message, next)
		},
		Condition: func(condition Condition[C]) {
			d.AddCondition(state, condition)
		},
		Decision: func(condition Condition[C], action Action[C, M]) {
			decision := Decision[C, M]{
				Condition: condition,
				Action:    action,
			}

			d.AddDecision(state, decision)
		},
		Fallback: func(fallback Action[C, M]) {
			d.AddFallback(state, fallback)
		},
	}
}

func (d *Stago[C, S, M]) State() S {
	return d.fsm.State()
}

func (d *Stago[C, S, M]) Context() C {
	return d.context
}

func (d *Stago[C, S, M]) Forward(message M) bool {
	var (
		state                       = d.State()
		actions, actions_exists     = d.actions[message]
		condition, condition_exists = d.conditions[state]
		decisions, decisions_exists = d.decisions[state]
		fallback, fallback_exists   = d.fallbacks[state]
	)

	if actions_exists {
		// Perform global actions
		for _, action := range actions {
			action(&d.context, message)
		}
	}

	if decisions_exists {
		// Perform conditional actions
		decided := false
		for _, decision := range decisions {
			if decision.Condition(d.context) {
				decision.Action(&d.context, message)
				decided = true
			}
		}
		if fallback_exists && !decided {
			// if neither decisions are made, perform fallback
			fallback(&d.context, message)
		}
	}

	if condition_exists {
		// Check if condition is true else abort transition
		if condition(d.context) {
			return false
		}
	}

	return d.fsm.Next(message)
}
