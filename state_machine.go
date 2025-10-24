/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import (
	"context"
	"errors"
	"io"
	"sync"
)

type (
	Transition[State comparable, Data any] struct {
		Previous State
		Next     *State
		Apply    []func(ctx context.Context, data Data) (Data, error)
	}

	StateMachine[State comparable, Data any] interface {
		Add(t *Transition[State, Data]) error
		Apply(ctx context.Context, state State, data Data) error
	}

	_stateMachine[State comparable, Data any] struct {
		states map[State]*Transition[State, Data]
		mux    sync.RWMutex
	}
)

func NewStateMachine[State comparable, Data any]() StateMachine[State, Data] {
	return &_stateMachine[State, Data]{
		states: make(map[State]*Transition[State, Data], 6),
	}
}

func (sm *_stateMachine[State, Data]) Add(t *Transition[State, Data]) error {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	if t == nil {
		return errors.New("transition cannot be nil")
	}

	if len(t.Apply) == 0 {
		return errors.New("transition must have at least one apply")
	}

	if _, ok := sm.states[t.Previous]; ok {
		return errors.New("transition already has a previous state")
	}

	if t.Next != nil {
		if t.Previous == *t.Next {
			return errors.New("transition has a identical previous and next state")
		}

		if tt, ok := sm.states[*t.Next]; ok && t.Previous == *tt.Next {
			return errors.New("transaction has a direct state loop")
		}
	}

	sm.states[t.Previous] = t
	return nil
}

func (sm *_stateMachine[State, Data]) Apply(ctx context.Context, state State, data Data) error {
	sm.mux.RLock()
	defer sm.mux.RUnlock()

	var err error
	for {
		t, ok := sm.states[state]
		if !ok {
			return nil
		}

		for _, apply := range t.Apply {
			e := Recovery(func() {
				data, err = apply(ctx, data)
			})
			if e != nil {
				return e
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return err
			}
		}

		newState := t.Next
		if newState == nil {
			return nil
		}
		state = *newState
	}
}
