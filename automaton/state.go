package automaton

type State interface {
	Enter()
	Exit()
	Next(any) State
}

type state[S any] struct {
	context  S
	on_enter OnEnter[S]
	on_exit  OnExit[S]
	on_next  OnNext[S]
}

type OnEnter[C any] func(ctx C)
type OnExit[C any] func(ctx C)
type OnNext[C any] func(ctx C, event any) State

type StateBuilder[S any] interface {
	OnEnter(OnEnter[S]) StateBuilder[S]
	OnExit(OnExit[S]) StateBuilder[S]
	OnNext(OnNext[S]) StateBuilder[S]
	Build() State
}

func NewStateBuilder[S any](ctx S) StateBuilder[S] {
	return &state[S]{
		context:  ctx,
		on_enter: nil,
		on_exit:  nil,
		on_next:  nil,
	}
}

func (sb *state[S]) OnEnter(callback OnEnter[S]) StateBuilder[S] {
	sb.on_enter = callback
	return sb
}

func (sb *state[S]) OnExit(callback OnExit[S]) StateBuilder[S] {
	sb.on_exit = callback
	return sb
}

func (sb *state[S]) OnNext(callback OnNext[S]) StateBuilder[S] {
	sb.on_next = callback
	return sb
}

func (sb *state[S]) Build() State {
	return sb
}

func (s *state[S]) Enter() {
	if s.on_enter != nil {
		s.on_enter(s.context)
	}
}

func (s *state[S]) Exit() {
	if s.on_exit != nil {
		s.on_exit(s.context)
	}
}

func (s *state[S]) Next(event any) State {
	if s.on_next != nil {
		return s.on_next(s.context, event)
	}
	return nil
}
