/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

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
	goDir := os.Getenv("GOROOT")

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
