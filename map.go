/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import (
	"slices"
)

func EachMap[K Comparable, V any](in map[K]V, call func(key K, value V)) {
	for _, k := range Keys[K, V](in) {
		call(k, in[k])
	}
}

func FilterMap[K comparable, V any](in map[K]V, filter func(key K, value V) bool) (out map[K]V) {
	out = make(map[K]V, len(in))
	for k, v := range in {
		if !filter(k, v) {
			continue
		}
		out[k] = v
	}
	return
}

func Keys[K Comparable, V any](in map[K]V) (out []K) {
	out = make([]K, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	slices.Sort(out)
	return
}

func Values[K Comparable, V any](in map[K]V) (out []V) {
	out = make([]V, 0, len(in))
	for _, k := range Keys[K, V](in) {
		out = append(out, in[k])
	}
	return
}

func JoinMap[K comparable, V any](in ...map[K]V) (out map[K]V) {
	size := 0
	for _, m := range in {
		size += len(m)
	}
	out = make(map[K]V, size)
	for _, m := range in {
		for k, v := range m {
			out[k] = v
		}
	}
	return
}

func FlipMap[K, V comparable](in map[K]V) (out map[V]K) {
	out = make(map[V]K, len(in))
	for k, v := range in {
		out[v] = k
	}
	return
}

func DivideMap[K Comparable, V any](in map[K]V) (keys []K, values []V) {
	keys = Keys[K, V](in)
	values = make([]V, 0, len(in))
	for _, k := range keys {
		values = append(values, in[k])
	}
	return
}

func CombineMap[K comparable, V any](keys []K, values []V) (out map[K]V) {
	out = make(map[K]V, len(keys))
	lenValues := len(values) - 1
	for i, k := range keys {
		var zero V
		if i > lenValues {
			out[k] = zero
			continue
		}
		out[k] = values[i]
	}
	return
}

func ReduceMap[K Comparable, V any](in map[K]V, call func(result, value V, key K) V) (out V) {
	for _, k := range Keys[K, V](in) {
		out = call(out, in[k], k)
	}
	return
}
