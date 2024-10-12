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

func TestUnit_CombineMap(t *testing.T) {
	out := do.CombineMap[int, string]([]int{1, 2, 3}, []string{"1", "2", "3", "4"})
	casecheck.Equal(t, map[int]string{1: "1", 2: "2", 3: "3"}, out)

	out = do.CombineMap[int, string]([]int{1, 2, 3}, []string{"1", "2"})
	casecheck.Equal(t, map[int]string{1: "1", 2: "2", 3: ""}, out)
}

func TestUnit_EachMap(t *testing.T) {
	result := "|"
	do.EachMap[string, int](map[string]int{"1": 1, "2": 2}, func(key string, value int) {
		result += fmt.Sprintf("%s:%d|", key, value)
	})
	casecheck.Equal(t, `|1:1|2:2|`, result)
}

func TestUnit_FilterMap(t *testing.T) {
	out := do.FilterMap[int, string](map[int]string{1: "1", 2: "2", 3: "3"}, func(key int, _ string) bool {
		return key%2 == 0
	})
	casecheck.Equal(t, map[int]string{2: "2"}, out)
}

func TestUnit_Keys(t *testing.T) {
	out := do.Keys[int, string](map[int]string{1: "1", 2: "2", 3: "3"})
	casecheck.Equal(t, []int{1, 2, 3}, out)
}

func TestUnit_Values(t *testing.T) {
	out := do.Values[int, string](map[int]string{1: "1", 2: "2", 3: "3"})
	casecheck.Equal(t, []string{"1", "2", "3"}, out)
}

func TestUnit_JoinMap(t *testing.T) {
	out := do.JoinMap[int, string](map[int]string{1: "1", 2: "2", 3: "3"}, map[int]string{1: "1", 2: "22", 4: "4"})
	casecheck.Equal(t, map[int]string{1: "1", 2: "22", 3: "3", 4: "4"}, out)
}

func TestUnit_FlipMap(t *testing.T) {
	out := do.FlipMap[int, string](map[int]string{1: "01", 2: "02", 3: "03"})
	casecheck.Equal(t, map[string]int{"01": 1, "02": 2, "03": 3}, out)
}

func TestUnit_DivideMap(t *testing.T) {
	keys, values := do.DivideMap[int, string](map[int]string{2: "02", 3: "03", 1: "01"})
	casecheck.Equal(t, []int{1, 2, 3}, keys)
	casecheck.Equal(t, []string{"01", "02", "03"}, values)
}

func TestUnit_ReduceMap(t *testing.T) {
	out := do.ReduceMap[int, string](map[int]string{2: "02", 1: "01", 3: "03"}, func(result, value string, key int) string {
		return result + value
	})
	casecheck.Equal(t, "010203", out)
}

func TestUnit_ToSlice(t *testing.T) {
	type A struct {
		K int
		V string
	}
	out := do.ToSlice[int, string, A](map[int]string{2: "02", 1: "01", 3: "03"}, func(value string, key int) A {
		return A{
			K: key,
			V: value,
		}
	})
	casecheck.Equal(t, []A{
		{K: 1, V: "01"},
		{K: 2, V: "02"},
		{K: 3, V: "03"},
	}, out)
}
