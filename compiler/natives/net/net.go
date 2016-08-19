// +build js

package net

import (
	"strings"
	"syscall"

	"github.com/bpowers/browsix-gopherjs/js"
)

func runtime_pollServerInit() {
}

func runtime_pollOpen(fd uintptr) (uintptr, int) {
	return 1, 0
}

func runtime_pollClose(ctx uintptr) {
}

func runtime_pollWait(ctx uintptr, mode int) int {
	return 0
}

func runtime_pollWaitCanceled(ctx uintptr, mode int) int {
	return 0
}

func runtime_pollReset(ctx uintptr, mode int) int {
	return 0
}

func runtime_pollSetDeadline(ctx uintptr, d int64, mode int) {
}

func runtime_pollUnblock(ctx uintptr) {
}

func byteIndex(s string, c byte) int {
	return js.InternalObject(s).Call("indexOf", js.Global.Get("String").Call("fromCharCode", c)).Int()
}

func listenBrowsix(net, laddr string) (Listener, error) {
	// FIXME: currently only support stream sockets
	family, sotype := syscall.AF_INET, syscall.SOCK_STREAM

	if parts := strings.SplitN(laddr, ":", 2); len(parts) == 2 && parts[0] == "localhost" {
		laddr = "127.0.0.1:" + parts[1]
	}

	addr, err := ResolveTCPAddr("tcp4", laddr)
	if err != nil {
		return nil, err
	}

	s, err := syscall.Socket(family, sotype, 0)
	if err != nil {
		return nil, err
	}

	fd, err := newFD(s, family, sotype, "tcp")
	if err != nil {
		return nil, err
	}

	if err = fd.listenStream(addr, 511); err != nil {
		return nil, err
	}

	//if err = syscall.Listen(s, 511); err != nil {
	//	return nil, err
	//}

	return &TCPListener{fd}, nil
}

func Listen(net, laddr string) (Listener, error) {
	return listenBrowsix(net, laddr)
}

func (d *Dialer) Dial(network, address string) (Conn, error) {
	family, sotype := syscall.AF_INET, syscall.SOCK_STREAM

	if parts := strings.SplitN(address, ":", 2); len(parts) == 2 && parts[0] == "localhost" {
		address = "127.0.0.1:" + parts[1]
	}

	raddr, err := ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	s, err := syscall.Socket(family, sotype, 0)
	if err != nil {
		return nil, err
	}

	fd, err := newFD(s, family, sotype, "tcp")
	if err != nil {
		return nil, err
	}

	if err := fd.dial(nil, raddr, noDeadline, noCancel); err != nil {
		fd.Close()
		return nil, err
	}

	return newTCPConn(fd), nil
}

func sysInit() {
}

func probeIPv4Stack() bool {
	return false
}

func probeIPv6Stack() (supportsIPv6, supportsIPv4map bool) {
	return false, false
}

func probeWindowsIPStack() (supportsVistaIP bool) {
	return false
}

func maxListenerBacklog() int {
	return syscall.SOMAXCONN
}
