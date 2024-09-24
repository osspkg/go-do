/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import "fmt"

func Try(try func(), catch func(err error), finally func()) {
	if try != nil {
		err := Recovery(try)
		if err != nil && catch != nil {
			//nolint:errcheck
			_ = Recovery(func() {
				catch(err)
			})
		}
	}

	if finally != nil {
		//nolint:errcheck
		_ = Recovery(finally)
	}
}

type (
	StepByStep[V any] struct {
		steps []func(V) (V, error)
	}
)

func NewStepByStep[V any]() *StepByStep[V] {
	return &StepByStep[V]{
		steps: make([]func(V) (V, error), 0),
	}
}

func (v *StepByStep[V]) Add(fn func(V) (V, error)) {
	v.steps = append(v.steps, fn)
}

func (v *StepByStep[V]) Exec(value V) (val V, err error) {
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
