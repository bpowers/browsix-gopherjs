// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signal

import (
//	"os"
//	"syscall"
)

var sigch = make(chan uint32)

// Defined by the runtime package.
func signal_disable(uint32) {
}

func signal_enable(uint32) {
}

func signal_ignore(uint32) {
}

func signal_recv() uint32 {
	return <-sigch
}
