package automaton

type Automaton[M any] struct {
	state State[M]
}

func New[M any](initial State[M]) *Automaton[M] {
	return &Automaton[M]{
		state: initial,
	}
}

func (a *Automaton[M]) Forward(message M) {
	nextState := a.state.Next(message)
	if nextState != nil {
		a.ChangeState(nextState)
	}
}

func (a *Automaton[M]) State() State[M] {
	return a.state
}

func (a *Automaton[M]) ChangeState(next State[M]) {
	a.state.Exit()
	a.state = next
	a.state.Enter()
}
