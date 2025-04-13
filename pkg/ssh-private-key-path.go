// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

type SSHPrivateKey string

func (f SSHPrivateKey) String() string {
	return string(f)
}
