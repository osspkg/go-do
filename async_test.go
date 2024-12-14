/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"go.osspkg.com/casecheck"
	"go.osspkg.com/do"
)

func TestUnit_Recovery(t *testing.T) {
	err := do.Recovery(func() {
		func() {
			panic(1)
		}()
	})
	casecheck.Error(t, err)
	casecheck.Contains(t, err.Error(), "panic=1 trace=./async_test.go:24 go.osspkg.com/do_test.TestUnit_Recovery.func1.1")
	casecheck.Equal(t, "1", errors.Unwrap(err).Error())
}

func TestUnit_Async(t *testing.T) {
	var err atomic.Value
	do.Async(func() {
		func() {
			panic(1)
		}()
	}, func(e error) {
		err.Store(e)
	})
	time.Sleep(100 * time.Millisecond)
	casecheck.Contains(t, err.Load().(error).Error(), "panic=1 trace=./async_test.go:36 go.osspkg.com/do_test.TestUnit_Async.func1.1")
	casecheck.Equal(t, "1", errors.Unwrap(err.Load().(error)).Error())
}

func TestUnit_AsyncGroup(t *testing.T) {
	errs := do.AsyncGroup(
		context.TODO(),
		func(ctx context.Context) error {
			panic(1)
		},
		func(ctx context.Context) error {
			panic(2)
		},
		func(ctx context.Context) error {
			return fmt.Errorf("good")
		},
	)
	var strErrs []string
	for _, err := range errs {
		strErrs = append(strErrs, err.Error())
	}
	sort.Strings(strErrs)
	casecheck.Equal(t, []string{"1", "2", "good"}, strErrs)
}
