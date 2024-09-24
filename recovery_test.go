/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do_test

import (
	"errors"
	"testing"

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
	casecheck.Contains(t, err.Error(), "panic=1 trace=./recovery_test.go:19 go.osspkg.com/do_test.TestUnit_Recovery.func1.1")
	casecheck.Equal(t, "1", errors.Unwrap(err).Error())
}
