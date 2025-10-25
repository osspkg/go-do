/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

func Each[T any](in []T, call func(value T, index int)) {
	for i, v := range in {
		call(v, i)
	}
}

func Convert[T, V any](in []T, call func(value T, index int) V) (out []V) {
	out = make([]V, len(in))
	for i, v := range in {
		out[i] = call(v, i)
	}
	return
}

func Join[T any](in ...[]T) (out []T) {
	size := 0
	for _, m := range in {
		size += len(m)
	}
	out = make([]T, 0, size)
	for _, s := range in {
		out = append(out, s...)
	}
	return
}

func Chunk[T any](in []T, count int) (out [][]T) {
	count = MinMax(count, 1, len(in))
	out = make([][]T, 0, len(in)/count+IfElse(len(in)%count == 0, 0, 1))

	chunk := make([]T, 0, 2)
	for i, v := range in {
		if i%count == 0 {
			if len(chunk) != 0 {
				out = append(out, chunk)
			}
			chunk = make([]T, 0, count)
		}
		chunk = append(chunk, v)
	}
	out = append(out, chunk)
	return
}

func ToMap[T comparable](in []T) (out map[T]struct{}) {
	out = make(map[T]struct{}, len(in))
	for _, v := range in {
		out[v] = struct{}{}
	}
	return
}

func Entries[T any, K comparable, V any](in []T, call func(T) (K, V)) (out map[K]V) {
	out = make(map[K]V, len(in))
	for _, val := range in {
		k, v := call(val)
		out[k] = v
	}
	return
}

func Reduce[T any](in []T, call func(result, value T, index int) T) (out T) {
	for i, v := range in {
		out = call(out, v, i)
	}
	return
}

func LastIndexOf[T comparable](in []T, element T) (index int, ok bool) {
	for i, v := range in {
		if v == element {
			index = i
			ok = true
		}
	}
	return
}

func IndexOf[T comparable](in []T, element T) (index int, ok bool) {
	for i, v := range in {
		if v == element {
			return i, true
		}
	}
	return 0, false
}

func Include[T comparable](in []T, element T) (ok bool) {
	_, ok = IndexOf[T](in, element)
	return
}

func Exclude[T comparable](in []T, elements ...T) (out []T) {
	out = make([]T, 0, len(in))
	elMap := ToMap[T](elements)
	for _, t := range in {
		if _, ok := elMap[t]; ok {
			continue
		}
		out = append(out, t)
	}
	return
}

func Pop[T any](in *[]T) (out T, ok bool) {
	if len(*in) == 0 {
		return
	}

	i := len(*in) - 1
	out, ok = (*in)[i], true
	*in = (*in)[:i]
	return
}

func Push[T any](in *[]T, elements ...T) {
	*in = append(*in, elements...)
}

func Shift[T any](in *[]T) (out T, ok bool) {
	if len(*in) == 0 {
		return
	}

	li := len(*in) - 1
	out, ok = (*in)[0], true
	if li == 0 {
		*in = (*in)[:0]
	} else {
		*in = (*in)[1:]
	}
	return
}

func Unshift[T any](in *[]T, elements ...T) {
	Splice[T](in, 0, 0, elements...)
}

func Reverse[T any](in []T) {
	j := len(in) - 1
	for i := 0; i < len(in)/2; i++ {
		in[i], in[j] = in[j], in[i]
		j--
	}
}

// Copy - start, end - element index (inclusive)
func Copy[T any](in []T, start, end int) (out []T) {
	start = MinMax(start, 0, len(in))
	end = MinMax(end+1, start+1, len(in))
	out = make([]T, end-start)
	copy(out[0:], in[start:end])
	return
}

func Splice[T any](in *[]T, start, deleteCount int, elements ...T) {
	lastPos := len(*in)
	start = MinMax(start, 0, lastPos)
	deleteCount = MinMax(deleteCount, 0, lastPos)
	deletePos := MinMax(deleteCount+start, 0, lastPos)

	*in = append((*in)[:start], (*in)[deletePos:]...)
	lastPos = len(*in)

	*in = append(*in, elements...)
	copy((*in)[start+len(elements):], (*in)[start:lastPos])
	copy((*in)[start:], elements[:])
}

func Filter[T any](in []T, filter func(value T, index int) bool) (out []T) {
	out = make([]T, 0, len(in))
	for i, v := range in {
		if filter(v, i) {
			out = append(out, v)
		}
	}
	return
}

func Treat[T any](in []T, prepare ...func(value T, index int) T) (out []T) {
	out = make([]T, 0, len(in))
	for i, v := range in {
		for _, fn := range prepare {
			v = fn(v, i)
		}
		out = append(out, v)
	}
	return
}

func TreatValue[T any](in []T, prepare ...func(T) T) (out []T) {
	out = make([]T, 0, len(in))
	for _, v := range in {
		for _, fn := range prepare {
			v = fn(v)
		}
		out = append(out, v)
	}
	return
}

func Diff[T comparable](arr1 []T, arr2 []T) (out []T) {
	out = make([]T, 0, len(arr1)+len(arr2))
	arr2Map := ToMap[T](arr2)
	for _, v := range arr1 {
		if _, ok := arr2Map[v]; ok {
			delete(arr2Map, v)
			continue
		}
		out = append(out, v)
	}
	for _, v := range arr2 {
		if _, ok := arr2Map[v]; ok {
			out = append(out, v)
		}
	}
	return
}

func Unique[T comparable](in []T) (out []T) {
	inMap := ToMap[T](in)
	out = make([]T, 0, len(inMap))
	for _, v := range in {
		if _, ok := inMap[v]; ok {
			out = append(out, v)
			delete(inMap, v)
		}
	}
	return
}
