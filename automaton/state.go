package automaton

type State[M any] interface {
	Enter()
	Exit()
	Next(M) State[M]
}

type state[S any, M any] struct {
	context  S
	on_enter OnEnter[S]
	on_exit  OnExit[S]
	on_next  OnNext[S, M]
}

type OnEnter[C any] func(ctx C)
type OnExit[C any] func(ctx C)
type OnNext[C any, M any] func(ctx C, event M) State[M]

type StateBuilder[S any, M any] interface {
	OnEnter(OnEnter[S]) StateBuilder[S, M]
	OnExit(OnExit[S]) StateBuilder[S, M]
	OnNext(OnNext[S, M]) StateBuilder[S, M]
	Build() State[M]
}

func NewStateBuilder[M any, S any](ctx S) StateBuilder[S, M] {
	return &state[S, M]{
		context:  ctx,
		on_enter: nil,
		on_exit:  nil,
		on_next:  nil,
	}
}

func (sb *state[S, M]) OnEnter(callback OnEnter[S]) StateBuilder[S, M] {
	sb.on_enter = callback
	return sb
}

func (sb *state[S, M]) OnExit(callback OnExit[S]) StateBuilder[S, M] {
	sb.on_exit = callback
	return sb
}

func (sb *state[S, M]) OnNext(callback OnNext[S, M]) StateBuilder[S, M] {
	sb.on_next = callback
	return sb
}

func (sb *state[S, M]) Build() State[M] {
	return sb
}

func (s *state[S, M]) Enter() {
	if s.on_enter != nil {
		s.on_enter(s.context)
	}
}

func (s *state[S, M]) Exit() {
	if s.on_exit != nil {
		s.on_exit(s.context)
	}
}

func (s *state[S, M]) Next(event M) State[M] {
	if s.on_next != nil {
		return s.on_next(s.context, event)
	}
	return nil
}
