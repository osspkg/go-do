/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"fmt"
	"testing"

	"go.osspkg.com/casecheck"

	"go.osspkg.com/do"
)

func TestUnit_StepByStep(t *testing.T) {
	sbs := do.NewStepByStep[string]()
	sbs.Add(func(s string) (string, error) {
		return s + "a", nil
	})
	sbs.Add(func(s string) (string, error) {
		return s + "b", nil
	})

	v, err := sbs.Exec("123")
	casecheck.NoError(t, err)
	casecheck.Equal(t, "123ab", v)

	sbs.Add(func(s string) (string, error) {
		switch s {
		case "000ab":
			return s, fmt.Errorf("1")
		case "111ab":
			panic("0")
		default:
			return s, nil
		}
	})

	_, err = sbs.Exec("000")
	casecheck.Error(t, err)
	casecheck.Equal(t, "fail on step #3: 1", err.Error())

	_, err = sbs.Exec("111")
	casecheck.Error(t, err)
	casecheck.Contains(t, err.Error(), "panic on step #3: panic=0 trace=./step_by_step_test.go")
	casecheck.Contains(t, err.Error(), "go.osspkg.com/do_test.TestUnit_StepByStep.func3")
}
