/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Words interface {
	Strings() []string
	Bytes() [][]byte
	UseDefaultBlock()
	SetBlock(detectors ...func(r rune) bool)
	UseDefaultDigital()
	SetDigital(detectors ...func(r rune) bool)
	UseDefaultSymbol()
	SetSymbol(detectors ...func(r rune) bool)
}

type _words struct {
	reader io.ReadSeeker
	block  []func(r rune) bool
	symbol []func(r rune) bool
	digit  []func(r rune) bool
}

func NewStringWords(s string) Words {
	w := &_words{
		reader: strings.NewReader(s),
	}
	w.UseDefaultBlock()
	w.UseDefaultDigital()
	w.UseDefaultSymbol()
	return w
}

func NewBytesWords(b []byte) Words {
	w := &_words{
		reader: bytes.NewReader(b),
	}
	w.UseDefaultBlock()
	w.UseDefaultDigital()
	w.UseDefaultSymbol()
	return w
}

func (v *_words) Strings() (out []string) {
	//nolint:errcheck
	v.reader.Seek(0, 0)
	scanner := bufio.NewScanner(v.reader)
	scanner.Split(v.scannerFunc)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	return
}

func (v *_words) Bytes() (out [][]byte) {
	//nolint:errcheck
	v.reader.Seek(0, 0)
	scanner := bufio.NewScanner(v.reader)
	scanner.Split(v.scannerFunc)
	for scanner.Scan() {
		out = append(out, scanner.Bytes())
	}
	return
}

func (v *_words) UseDefaultBlock() {
	v.SetBlock(unicode.IsDigit, unicode.IsLetter)
}

func (v *_words) SetBlock(detectors ...func(r rune) bool) {
	v.block = detectors
}

func (v *_words) isBlock(r rune) bool {
	for _, fn := range v.block {
		if fn(r) {
			return true
		}
	}
	return false
}

func (v *_words) UseDefaultDigital() {
	v.SetDigital(unicode.IsDigit)
}

func (v *_words) SetDigital(detectors ...func(r rune) bool) {
	v.digit = detectors
}

func (v *_words) isDigital(r rune) bool {
	for _, fn := range v.digit {
		if fn(r) {
			return true
		}
	}
	return false
}

func (v *_words) UseDefaultSymbol() {
	v.SetSymbol(unicode.IsPunct, unicode.IsSymbol)
}

func (v *_words) SetSymbol(detectors ...func(r rune) bool) {
	v.symbol = detectors
}

func (v *_words) isSymbol(r rune) bool {
	for _, fn := range v.symbol {
		if fn(r) {
			return true
		}
	}
	return false
}

func (v *_words) scannerFunc(data []byte, atEOF bool) (int, []byte, error) {
	var (
		start int
	)

	for index := 0; index < len(data); {
		cR, cW := utf8.DecodeRune(data[index:])
		nR, _ := utf8.DecodeRune(data[index+cW:])
		index += cW

		switch true {
		case v.isBlock(cR):
			if !v.isBlock(nR) {
				return index, data[start:index], nil
			}
		case v.isDigital(cR):
			if !v.isDigital(nR) {
				return index, data[start:index], nil
			}
		case v.isSymbol(cR):
			if !v.isSymbol(nR) {
				return index, data[start:index], nil
			}
		default:
			start = index
		}
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}

func WordStrings(s string) []string {
	w := NewStringWords(s)
	w.UseDefaultBlock()
	w.UseDefaultSymbol()
	return w.Strings()
}

func WordBytes(b []byte) [][]byte {
	w := NewBytesWords(b)
	w.UseDefaultBlock()
	w.UseDefaultSymbol()
	return w.Bytes()
}
