/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import "time"

func MinMax[T Comparable](value, minimum, maximum T) T {
	if minimum > value {
		return minimum
	}
	if maximum < value {
		return maximum
	}
	return value
}

func MinMaxTime(value, minimum, maximum time.Time) time.Time {
	if value.Before(minimum) {
		return minimum
	}
	if value.After(maximum) {
		return maximum
	}
	return value
}

func Range[T Summable](from, to, step T) (out []T) {
	out = make([]T, 0, 2)
	for i := from; i <= to; i += step {
		out = append(out, i)
	}
	return
}

func Sum[T Comparable](elements ...T) (out T) {
	for _, element := range elements {
		out += element
	}
	return
}

func Average[T Summable](elements ...T) (out T) {
	for _, element := range elements {
		out += element
	}
	out = out / T(len(elements))
	return
}
