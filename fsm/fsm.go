package fsm

// Transition Table.
type TT[S, M comparable] map[S]TR[S, M]

func (tt TT[S, M]) InitState(s S) {
	tt[s] = make(TR[S, M])
}

func (tt TT[S, M]) AddTransition(s S, m M, ns S) {
	if _, ok := tt[s]; !ok {
		tt.InitState(s)
	}

	tt[s][m] = ns
}

// Transition Row.
type TR[S, M comparable] map[M]S

// Finite-state machine.
//
//	S - State
//	M - Message
type FSM[S, M comparable] struct {
	state S
	table TT[S, M]
}

// New allocates FSM.
func New[S, M comparable](init S, table TT[S, M]) *FSM[S, M] {
	return &FSM[S, M]{
		state: init,
		table: table,
	}
}

// Next processes the message and changes the state according to provided transition table.
func (f *FSM[S, M]) Next(message M) bool {
	next, ok := f.table[f.state][message]
	if ok {
		f.state = next
	}

	return ok
}

func (f *FSM[S, M]) Table() TT[S, M] {
	return f.table
}

func (f *FSM[S, M]) Reset(state S) {
	f.state = state
}

// State returns the FSM's state.
func (f *FSM[S, M]) State() S {
	return f.state
}
