/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"fmt"
	"testing"

	"go.osspkg.com/casecheck"
	"go.osspkg.com/do"
)

func TestUnit_Each(t *testing.T) {
	result := "|"
	do.Each[string]([]string{"1", "2"}, func(key string, index int) {
		result += fmt.Sprintf("(%d)%s|", index, key)
	})
	casecheck.Equal(t, `|(0)1|(1)2|`, result)
}

func TestUnit_Join(t *testing.T) {
	out := do.Join[string]([]string{"1", "2"}, []string{"5", "6"}, []string{"0", "2"})
	casecheck.Equal(t, []string{"1", "2", "5", "6", "0", "2"}, out)
}

func TestUnit_Chunk(t *testing.T) {
	out := do.Chunk[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1}, 3)
	casecheck.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {0, 1}}, out)

	out = do.Chunk[int]([]int{1, 2, 3}, 3)
	casecheck.Equal(t, [][]int{{1, 2, 3}}, out)

	out = do.Chunk[int]([]int{1, 2, 3}, 0)
	casecheck.Equal(t, [][]int{{1}, {2}, {3}}, out)
}

func TestUnit_ToMap(t *testing.T) {
	out := do.ToMap[string]([]string{"1", "2"})
	casecheck.Equal(t, map[string]struct{}{"1": {}, "2": {}}, out)
}

func TestUnit_Reduce(t *testing.T) {
	out := do.Reduce[string]([]string{"1", "2"}, func(result, value string, index int) string {
		return result + value
	})
	casecheck.Equal(t, "12", out)
}

func TestUnit_LastIndexOf(t *testing.T) {
	out, ok := do.LastIndexOf[int]([]int{1, 2, 3, 3, 4}, 3)
	casecheck.True(t, ok)
	casecheck.Equal(t, 3, out)
}

func TestUnit_IndexOf(t *testing.T) {
	out, ok := do.IndexOf[int]([]int{1, 2, 3, 3, 4}, 3)
	casecheck.True(t, ok)
	casecheck.Equal(t, 2, out)

	out, ok = do.IndexOf[int]([]int{1, 2, 3, 3, 4}, 10)
	casecheck.False(t, ok)
}

func TestUnit_Include(t *testing.T) {
	ok := do.Include[int]([]int{1, 2, 3, 3, 4}, 3)
	casecheck.True(t, ok)

	ok = do.Include[int]([]int{1, 2, 3, 3, 4}, 10)
	casecheck.False(t, ok)
}

func TestUnit_Exclude(t *testing.T) {
	out := do.Exclude[int]([]int{1, 2, 3, 3, 4}, 3)
	casecheck.Equal(t, []int{1, 2, 4}, out)
}

func TestUnit_Pop(t *testing.T) {
	data := []int{1, 2}

	out, ok := do.Pop[int](&data)
	casecheck.True(t, ok)
	casecheck.Equal(t, 2, out)

	out, ok = do.Pop[int](&data)
	casecheck.True(t, ok)
	casecheck.Equal(t, 1, out)

	out, ok = do.Pop[int](&data)
	casecheck.False(t, ok)
}

func TestUnit_Push(t *testing.T) {
	data := []int{1, 2}

	do.Push[int](&data, 5, 6, 7)
	casecheck.Equal(t, []int{1, 2, 5, 6, 7}, data)
}

func TestUnit_Shift(t *testing.T) {
	data := []int{1, 2}

	out, ok := do.Shift[int](&data)
	casecheck.True(t, ok)
	casecheck.Equal(t, 1, out)

	out, ok = do.Shift[int](&data)
	casecheck.True(t, ok)
	casecheck.Equal(t, 2, out)

	out, ok = do.Shift[int](&data)
	casecheck.False(t, ok)
	casecheck.Equal(t, []int{}, data)
}

func TestUnit_Unshift(t *testing.T) {
	data := []int{1, 2}
	do.Unshift[int](&data, 10, 12, 11)
	casecheck.Equal(t, []int{10, 12, 11, 1, 2}, data)
}

func TestUnit_Reverse(t *testing.T) {
	data := []int{1, 2, 10, 12, 11}
	do.Reverse(data)
	casecheck.Equal(t, []int{11, 12, 10, 2, 1}, data)
}

func TestUnit_Splice(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	do.Splice[int](&data, 2, 3, 111, 112, 113)
	casecheck.Equal(t, []int{0, 1, 111, 112, 113, 5, 6, 7, 8, 9}, data)

	data = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	do.Splice[int](&data, 2, 10, 111, 112, 113)
	casecheck.Equal(t, []int{0, 1, 111, 112, 113}, data)

	data = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	do.Splice[int](&data, 100, 10, 111, 112, 113)
	casecheck.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 111, 112, 113}, data)

	data = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	do.Splice[int](&data, 0, 10, 111, 112, 113)
	casecheck.Equal(t, []int{111, 112, 113}, data)

	data = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	do.Splice[int](&data, 0, 0, 111, 112, 113)
	casecheck.Equal(t, []int{111, 112, 113, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, data)

	data = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	do.Splice[int](&data, -1, -1, 111, 112, 113)
	casecheck.Equal(t, []int{111, 112, 113, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, data)
}

func TestUnit_Copy(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	out := do.Copy[int](data, 2, 5)
	casecheck.Equal(t, []int{2, 3, 4, 5}, out)

	out = do.Copy[int](data, 2, 10)
	casecheck.Equal(t, []int{2, 3, 4, 5, 6, 7, 8, 9}, out)

	out = do.Copy[int](data, 2, 0)
	casecheck.Equal(t, []int{2}, out)
}

func TestUnit_Filter(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	out := do.Filter[int](data, func(value int, _ int) bool {
		return value%2 == 0
	})
	casecheck.Equal(t, []int{0, 2, 4, 6, 8}, out)
}

func TestUnit_Diff(t *testing.T) {
	data1 := []int{0, 1, 2, 3, 4}
	data2 := []int{3, 4, 5, 6, 7, 8, 9}

	out := do.Diff[int](data1, data2)
	casecheck.Equal(t, []int{0, 1, 2, 5, 6, 7, 8, 9}, out)
}

func TestUnit_Uniq(t *testing.T) {
	data := []int{0, 1, 5, 2, 1, 5, 2, 8, 0, 2}

	out := do.Unique[int](data)
	casecheck.Equal(t, []int{0, 1, 5, 2, 8}, out)
}

func TestUnit_Entries(t *testing.T) {
	data := []int{5}

	out := do.Entries[int, int, int](data, func(i int) (int, int) {
		return i, i * 2
	})
	casecheck.Equal(t, map[int]int{5: 10}, out)
}
