package automaton

type Automaton struct {
	state State
}

func New(initial State) *Automaton {
	return &Automaton{
		state: initial,
	}
}

func (a *Automaton) Forward(event any) {
	next_state := a.state.Next(event)
	if next_state != nil {
		a.ChangeState(next_state)
	}
}

func (a *Automaton) State() State {
	return a.state
}

func (a *Automaton) ChangeState(next State) {
	a.state.Exit()
	a.state = next
	a.state.Enter()
}
