/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

type _doIf[T any] struct {
	value T
	done  bool
}

type DoIf[T any] interface {
	ElseIf(compare bool, value T) DoIf[T]
	Else(value T) T
	ElseIfFunc(compare bool, value func() T) DoIf[T]
	ElseFunc(value func() T) T
}

func If[T any](compare bool, value T) DoIf[T] {
	result := &_doIf[T]{}

	if compare {
		result.done = true
		result.value = value
	}

	return result
}

func IfFunc[T any](compare bool, value func() T) DoIf[T] {
	result := &_doIf[T]{}

	if compare {
		result.done = true
		result.value = value()
	}

	return result
}

func (v *_doIf[T]) ElseIf(compare bool, value T) DoIf[T] {
	if !v.done && compare {
		v.done = true
		v.value = value
	}
	return v
}

func (v *_doIf[T]) ElseIfFunc(compare bool, value func() T) DoIf[T] {
	if !v.done && compare {
		v.done = true
		v.value = value()
	}
	return v
}

func (v *_doIf[T]) Else(value T) T {
	if v.done {
		return v.value
	}
	return value
}

func (v *_doIf[T]) ElseFunc(value func() T) T {
	if v.done {
		return v.value
	}
	return value()
}

func IfElse[T any](compare bool, thenV, elseV T) T {
	if compare {
		return thenV
	}
	return elseV
}

func IfElseFunc[T any](compare bool, thenV, elseV func() T) T {
	if compare {
		return thenV()
	}
	return elseV()
}
