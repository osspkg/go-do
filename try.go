/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package do

func Try(try func(), catch func(err error), finally func()) {
	if try != nil {
		err := Recovery(try)
		if err != nil && catch != nil {
			//nolint:errcheck
			_ = Recovery(func() {
				catch(err)
			})
		}
	}

	if finally != nil {
		//nolint:errcheck
		_ = Recovery(finally)
	}
}
