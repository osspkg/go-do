/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"testing"
	"time"

	"go.osspkg.com/casecheck"

	"go.osspkg.com/do"
)

func TestUnit_MinMax(t *testing.T) {
	casecheck.Equal(t, 5, do.MinMax[int](5, 0, 10))
	casecheck.Equal(t, 0, do.MinMax[int](-5, 0, 10))
	casecheck.Equal(t, 10, do.MinMax[int](20, 0, 10))
}

func TestUnit_Sum(t *testing.T) {
	casecheck.Equal(t, 15, do.Sum[int](5, 0, 10))
	casecheck.Equal(t, 11.5, do.Sum[float64](1.5, 0, 10))
	casecheck.Equal(t, "aabbcc", do.Sum[string]("aa", "bb", "cc"))
}

func TestUnit_Average(t *testing.T) {
	casecheck.Equal(t, 5, do.Average[int](5, 0, 10))
	casecheck.Equal(t, 3.75, do.Average[float64](1.5, 2.5, 1, 10))
}

func TestUnit_Range(t *testing.T) {
	casecheck.Equal(t, []float64{1, 1.5, 2, 2.5, 3}, do.Range[float64](1, 3, 0.5))
	casecheck.Equal(t, []int64{1, 2, 3}, do.Range[int64](1, 3, 1))
	casecheck.Equal(t, []byte{'a', 'b', 'c'}, do.Range[byte]('a', 'c', 1))
}

func TestUnit_MinMaxTime(t *testing.T) {
	mi := time.Now()
	ma := time.Now().Add(1 * time.Hour)

	v := mi.Add(-10 * time.Second)
	casecheck.Equal(t, mi, do.MinMaxTime(v, mi, ma))
	v = mi.Add(10 * time.Second)
	casecheck.Equal(t, v, do.MinMaxTime(v, mi, ma))
	v = ma.Add(10 * time.Second)
	casecheck.Equal(t, ma, do.MinMaxTime(v, mi, ma))
}
