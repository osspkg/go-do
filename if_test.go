/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"reflect"
	"testing"

	"go.osspkg.com/do"
)

func TestUnit_If(t *testing.T) {
	tests := []struct {
		name   string
		result any
		want   any
	}{
		{
			name:   "Case1",
			result: do.If[int](true, 1).Else(0),
			want:   1,
		},
		{
			name:   "Case2",
			result: do.If[int](false, 1).Else(0),
			want:   0,
		},
		{
			name:   "Case3",
			result: do.If[int](false, 1).ElseIf(true, 2).Else(0),
			want:   2,
		},
		{
			name:   "Case4",
			result: do.If[int](false, 1).ElseIf(false, 2).Else(0),
			want:   0,
		},
		{
			name:   "Case5",
			result: do.IfElse[int](false, 1, 0),
			want:   0,
		},
		{
			name:   "Case5",
			result: do.IfElse[int](true, 1, 0),
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.result, tt.want) {
				t.Errorf("If() = %v, want %v", tt.result, tt.want)
			}
		})
	}
}

func TestUnit_IfFunc(t *testing.T) {
	tests := []struct {
		name   string
		result any
		want   any
	}{
		{
			name:   "Case1",
			result: do.IfFunc[int](true, func() int { return 1 }).ElseFunc(func() int { return 0 }),
			want:   1,
		},
		{
			name:   "Case2",
			result: do.IfFunc[int](false, func() int { return 1 }).ElseFunc(func() int { return 0 }),
			want:   0,
		},
		{
			name:   "Case3",
			result: do.IfFunc[int](false, func() int { return 1 }).ElseIfFunc(true, func() int { return 2 }).ElseFunc(func() int { return 0 }),
			want:   2,
		},
		{
			name:   "Case4",
			result: do.IfFunc[int](false, func() int { return 1 }).ElseIfFunc(false, func() int { return 2 }).ElseFunc(func() int { return 0 }),
			want:   0,
		},
		{
			name:   "Case5",
			result: do.IfElseFunc[int](false, func() int { return 1 }, func() int { return 0 }),
			want:   0,
		},
		{
			name:   "Case5",
			result: do.IfElseFunc[int](true, func() int { return 1 }, func() int { return 0 }),
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.result, tt.want) {
				t.Errorf("If() = %v, want %v", tt.result, tt.want)
			}
		})
	}
}
