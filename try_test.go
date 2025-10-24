/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
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
