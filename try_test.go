/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"errors"
	"fmt"
	"testing"

	"go.osspkg.com/casecheck"
	"go.osspkg.com/do"
)

func TestUnit_Try(t *testing.T) {
	var errs error
	do.Try(nil, func(err error) {
		errs = err
	}, nil)
	casecheck.NoError(t, errs)

	errs = nil
	do.Try(func() {
		panic(1)
	}, func(err error) {
		errs = fmt.Errorf("%w+catch", errors.Unwrap(err))
	}, func() {
		errs = fmt.Errorf("%w+finally", errs)
	})
	casecheck.Equal(t, `1+catch+finally`, errs.Error())
}

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
	casecheck.Equal(t, "panic on step #3: panic=0 trace=./try_test.go:53 go.osspkg.com/do_test.TestUnit_StepByStep.func3", err.Error())
}
