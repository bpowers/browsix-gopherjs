// +build js

package runtime

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
