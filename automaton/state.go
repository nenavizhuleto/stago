package automaton

type State[M any] interface {
	Enter()
	Exit()
	Next(M) State[M]
}

type state[S any, M any] struct {
	context S
	onEnter OnEnter[S]
	onExit  OnExit[S]
	onNext  OnNext[S, M]
}

type (
	OnEnter[C any]       func(ctx C)
	OnExit[C any]        func(ctx C)
	OnNext[C any, M any] func(ctx C, event M) State[M]
)

type StateBuilder[S any, M any] interface {
	OnEnter(OnEnter[S]) StateBuilder[S, M]
	OnExit(OnExit[S]) StateBuilder[S, M]
	OnNext(OnNext[S, M]) StateBuilder[S, M]
	Build() State[M]
}

func NewStateBuilder[M any, S any](ctx S) StateBuilder[S, M] {
	return &state[S, M]{
		context: ctx,
		onEnter: nil,
		onExit:  nil,
		onNext:  nil,
	}
}

func (s *state[S, M]) OnEnter(callback OnEnter[S]) StateBuilder[S, M] {
	s.onEnter = callback

	return s
}

func (s *state[S, M]) OnExit(callback OnExit[S]) StateBuilder[S, M] {
	s.onExit = callback

	return s
}

func (s *state[S, M]) OnNext(callback OnNext[S, M]) StateBuilder[S, M] {
	s.onNext = callback

	return s
}

func (s *state[S, M]) Build() State[M] {
	return s
}

func (s *state[S, M]) Enter() {
	if s.onEnter != nil {
		s.onEnter(s.context)
	}
}

func (s *state[S, M]) Exit() {
	if s.onExit != nil {
		s.onExit(s.context)
	}
}

func (s *state[S, M]) Next(event M) State[M] {
	if s.onNext != nil {
		return s.onNext(s.context, event)
	}

	return nil
}
