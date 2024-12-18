/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func Recovery(call func()) (err error) {
	defer func() {
		if val := recover(); val != nil {
			err = fmt.Errorf(
				"panic=%w trace=%s",
				fmt.Errorf("%+v", val),
				Trace(
					4,
					1,
					"go.osspkg.com/do.Recovery",
				),
			)
		}
	}()

	call()
	return
}

var bufPool = sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) }}

func Trace(skipLines, countLines int, skipFunc ...string) string {
	list := make([]uintptr, countLines+1)

	//nolint:errcheck
	execFile, _ := os.Executable()
	execDir := filepath.Dir(execFile)
	//nolint:errcheck
	workDir, _ := os.Getwd()
	goDir := runtime.GOROOT()

	n := runtime.Callers(skipLines, list)
	frame := runtime.CallersFrames(list[:n])

	//nolint:errcheck
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	var line int
	for {
		v, ok := frame.Next()
		if !ok {
			break
		}

		doWrite := true
		for _, s := range skipFunc {
			if strings.Contains(v.Function, s) {
				doWrite = false
			}
		}

		if doWrite {
			line++
			filePath := strings.TrimPrefix(v.File, workDir)
			filePath = strings.TrimPrefix(filePath, execDir)
			filePath = strings.TrimPrefix(filePath, goDir)
			fmt.Fprintf(buf, "%s.%s:%d %v", IfElse(line == 1, "", "\n"), filePath, v.Line, v.Function)
		}

	}
	return buf.String()
}
