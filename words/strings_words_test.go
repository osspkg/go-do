/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package words_test

import (
	"testing"
	"unicode"

	"go.osspkg.com/casecheck"
	"go.osspkg.com/do/words"
)

func TestUnit_WordStrings(t *testing.T) {
	data := `
a11
Привет!
Hello World {user name }
123 * 123 =111
`
	out := words.Strings(data)
	casecheck.Equal(t, []string{
		"a11", "Привет", "!", "Hello", "World", "{", "user", "name",
		"}", "123", "*", "123", "=", "111",
	}, out)
}

func TestUnit_WordCustom(t *testing.T) {
	data := `
a11
Привет!
Hello World {user name }
123 * 123 =111
`
	w := words.NewString(data)
	w.SetBlock(unicode.IsLetter)
	out := w.Strings()
	casecheck.Equal(t, []string{
		"a", "11", "Привет", "!", "Hello", "World", "{", "user", "name",
		"}", "123", "*", "123", "=", "111",
	}, out)
}

func TestUnit_WordBytes(t *testing.T) {
	data := `
a11
Привет!
Hello World {user name }
123 * 123 =111
`
	out := words.Bytes([]byte(data))
	casecheck.Equal(t, [][]byte{
		[]byte("a11"),
		[]byte("Привет"),
		[]byte("!"),
		[]byte("Hello"),
		[]byte("World"),
		[]byte("{"),
		[]byte("user"),
		[]byte("name"),
		[]byte("}"),
		[]byte("123"),
		[]byte("*"),
		[]byte("123"),
		[]byte("="),
		[]byte("111"),
	}, out)
}
