// +build js

package pprof

import (
	"io"
	"sync"
)

type Profile struct {
	name  string
	mu    sync.Mutex
	m     map[interface{}][]uintptr
	count func() int
	write func(io.Writer, int) error
}

func Profiles() []*Profile {
	return []*Profile{} // nil
}

func (p *Profile) WriteTo(w io.Writer, debug int) error {
	return nil
}

func StartCPUProfile(w io.Writer) error {
	return nil
}

func StopCPUProfile() {
}

func WriteHeapProfile(w io.Writer) error {
	return nil
}

func Lookup(name string) *Profile {
	return nil
}
