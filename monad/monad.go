/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package monad

type Result[T any] struct {
	value T
	err   error
}

func Some[T any](arg T) Result[T] {
	return Result[T]{value: arg}
}

func Bind[T, R any](r Result[T], next func(T) (R, error)) Result[R] {
	if r.err != nil {
		return Result[R]{err: r.err}
	}

	v, err := next(r.value)
	if err != nil {
		return Result[R]{err: err}
	}

	return Result[R]{value: v}
}

func (r Result[T]) Return() (T, error) {
	return r.value, r.err
}
