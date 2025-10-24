/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

func Async(callFunc func(), errFunc func(err error)) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				if errFunc != nil {
					err := fmt.Errorf(
						"panic=%w trace=%s",
						fmt.Errorf("%+v", e),
						Trace(4, 1, "go.osspkg.com/do.Async"))
					errFunc(err)
				}
			}
		}()
		if callFunc != nil {
			callFunc()
		}
	}()
}

func AsyncGroup(ctx context.Context, callFuncs ...func(ctx context.Context) error) []error {
	var wg sync.WaitGroup
	errC := make(chan error, len(callFuncs))

	wg.Add(len(callFuncs))
	for _, callFunc := range callFuncs {
		callFunc := callFunc
		go func() {
			var err error
			defer func() {
				if e := recover(); e != nil {
					err = errors.Join(err, fmt.Errorf("%+v", e))
				}
				if err != nil {
					errC <- err
				}
				wg.Done()
			}()
			err = callFunc(ctx)
		}()
	}

	wg.Wait()
	close(errC)

	errGrp := make([]error, 0, len(callFuncs))
	for err := range errC {
		errGrp = append(errGrp, err)
	}

	return errGrp
}
