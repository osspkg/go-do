/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import "fmt"

type (
	StepByStep[V any] interface {
		Add(fn func(V) (V, error))
		Exec(value V) (val V, err error)
	}
	_stepByStep[V any] struct {
		steps []func(V) (V, error)
	}
)

func NewStepByStep[V any]() StepByStep[V] {
	return &_stepByStep[V]{
		steps: make([]func(V) (V, error), 0),
	}
}

func (v *_stepByStep[V]) Add(fn func(V) (V, error)) {
	v.steps = append(v.steps, fn)
}

func (v *_stepByStep[V]) Exec(value V) (val V, err error) {
	val = value
	for i, step := range v.steps {
		e := Recovery(func() {
			val, err = step(val)
		})
		if e != nil {
			err = fmt.Errorf("panic on step #%d: %w", i+1, e)
			return
		}
		if err != nil {
			err = fmt.Errorf("fail on step #%d: %w", i+1, err)
			return
		}
	}
	return
}
