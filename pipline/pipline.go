package pipline

import (
	"context"
	"errors"
)

type Pipe[S comparable, T any] struct {
	Current S
	Next    S
	Handler func(ctx context.Context, arg T) (T, error)
}

type Pipline[S comparable, T any] struct {
	pips map[S]*Pipe[S, T]
}

func New[S comparable, T any]() *Pipline[S, T] {
	return &Pipline[S, T]{
		pips: make(map[S]*Pipe[S, T]),
	}
}

func (p *Pipline[S, T]) Set(arg Pipe[S, T]) error {
	if arg.Handler == nil {
		return errors.New("nil handler")
	}
	if _, ok := p.pips[arg.Current]; ok {
		return errors.New("current state is already set")
	}
	if arg.Current == arg.Next {
		return errors.New("current and next state is identical")
	}

	p.pips[arg.Current] = &arg

	return nil
}

func (p *Pipline[S, T]) Do(ctx context.Context, state S, arg T) (S, T, error) {
	var err error

	for {
		select {
		case <-ctx.Done():
			return state, arg, ctx.Err()
		default:
		}

		ps, ok := p.pips[state]
		if !ok {
			return state, arg, nil
		}

		arg, err = ps.Handler(ctx, arg)
		if err != nil {
			return state, arg, err
		}

		state = ps.Next
	}
}
