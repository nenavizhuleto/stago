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
	actions    map[S]Action[C, M]
	fallbacks  map[S]Action[C, M]
}

func New[C any, S, M comparable](ctx C, state S) *Stago[C, S, M] {
	s := &Stago[C, S, M]{
		context:    ctx,
		fsm:        fsm.New(state, fsm.TT[S, M]{}),
		decisions:  make(map[S][]Decision[C, M]),
		conditions: make(map[S]Condition[C]),
		actions:    make(map[S]Action[C, M]),
		fallbacks:  make(map[S]Action[C, M]),
	}

	return s
}

type State[C any, S, M comparable] struct {
	v S
	s *Stago[C, S, M]
}

func (s State[C, S, M]) Transition(message M, next S) {
	s.s.fsm.Table().AddTransition(s.v, message, next)
}

func (s State[C, S, M]) Condition(condition Condition[C]) {
	s.s.AddCondition(s.v, condition)
}

func (s State[C, S, M]) Decision(condition Condition[C], action Action[C, M]) {
	decision := Decision[C, M]{
		Condition: condition,
		Action:    action,
	}

	s.s.AddDecision(s.v, decision)
}

func (s State[C, S, M]) Action(action Action[C, M]) {
	s.s.AddAction(s.v, action)
}

func (s State[C, S, M]) Fallback(fallback Action[C, M]) {
	s.s.AddFallback(s.v, fallback)
}

func (d *Stago[C, S, M]) AddDecision(state S, decision Decision[C, M]) {
	decisions := d.decisions[state]
	d.decisions[state] = append(decisions, decision)
}

func (d *Stago[C, S, M]) AddAction(state S, action Action[C, M]) {
	d.actions[state] = action
}

func (d *Stago[C, S, M]) AddFallback(state S, fallback Action[C, M]) {
	d.fallbacks[state] = fallback
}

func (d *Stago[C, S, M]) AddCondition(state S, condition Condition[C]) {
	d.conditions[state] = condition
}

func (d *Stago[C, S, M]) NewState(state S) State[C, S, M] {
	return State[C, S, M]{
		v: state,
		s: d,
	}
}

func (d *Stago[C, S, M]) State() S {
	return d.fsm.State()
}

func (d *Stago[C, S, M]) Context() C {
	return d.context
}

func (d *Stago[C, S, M]) ResetContext(ctx C) {
	d.context = ctx
}

func (d *Stago[C, S, M]) ResetState(state S) {
	d.fsm.Reset(state)
}

func (d *Stago[C, S, M]) Reset(ctx C, state S) {
	d.ResetContext(ctx)
	d.ResetState(state)
}

func (d *Stago[C, S, M]) Forward(message M) bool {
	var (
		state                       = d.State()
		action, action_exists       = d.actions[state]
		condition, condition_exists = d.conditions[state]
		decisions, decisions_exists = d.decisions[state]
		fallback, fallback_exists   = d.fallbacks[state]
	)

	if action_exists {
		// Perform actions
		action(&d.context, message)
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
		if !condition(d.context) {
			return false
		}
	}

	return d.fsm.Next(message)
}
