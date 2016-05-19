// +build js

package os

import (
	"github.com/bpowers/browsix-gopherjs/js"
)

func runtime_args() []string { // not called on Windows
	return Args
}

func init() {
	if process := js.Global.Get("process"); process != js.Undefined {
		argv := process.Get("argv")
		if argv == nil {
			ch := make(chan *js.Object, 0)
			process.Call("once", "ready", func (r *js.Object) {
				ch <- r
			})
			_ = <-ch
			argv = process.Get("argv")
		}
		Args = make([]string, argv.Length()-1)
		for i := 0; i < argv.Length()-1; i++ {
			Args[i] = argv.Index(i + 1).String()
		}
	}
	if len(Args) == 0 {
		Args = []string{"?"}
	}
}

func runtime_beforeExit() {}
