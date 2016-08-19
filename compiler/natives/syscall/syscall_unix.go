// +build js,!windows

package syscall

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/bpowers/browsix-gopherjs/js"
)

type sysret struct {
	r1, r2 uintptr
	err    Errno
}

var chans = sync.Pool{New: func() interface{} { return make(chan sysret) }}

func runtime_envs() []string {
	process := js.Global.Get("process")
	if process == js.Undefined {
		return nil
	}

	jsEnv := process.Get("env")
	if jsEnv == nil {
		ch := make(chan *js.Object, 0)
		process.Call("once", "ready", func(r *js.Object) {
			ch <- r
		})
		_ = <-ch
		jsEnv = process.Get("env")
	}

	envkeys := js.Global.Get("Object").Call("keys", jsEnv)
	envs := make([]string, envkeys.Length())
	for i := 0; i < envkeys.Length(); i++ {
		key := envkeys.Index(i).String()
		envs[i] = key + "=" + jsEnv.Get(key).String()
	}
	return envs
}

func setenv_c(k, v string) {
	process := js.Global.Get("process")
	if process != js.Undefined {
		process.Get("env").Set(k, v)
	}
}

var syscallModule *js.Object
var alreadyTriedToLoad = false
var minusOne = -1

func syscall(name string) *js.Object {
	defer func() {
		recover()
		// return nil if recovered
	}()
	if syscallModule == nil {
		if alreadyTriedToLoad {
			return nil
		}
		alreadyTriedToLoad = true
		require := js.Global.Get("require")
		if require == js.Undefined {
			panic("")
		}
		syscallModule = require.Invoke("syscall")
	}
	return syscallModule.Get(name)
}

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	if f := syscall("Syscall"); f != nil {
		ch := chans.Get().(chan sysret)
		f.Invoke(func(r *js.Object) {
			ch <- sysret{uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())}
		}, trap, a1, a2, a3)
		result := <-ch
		chans.Put(ch)
		return result.r1, result.r2, result.err
		// uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())
	}
	if trap == SYS_EXIT {
		runtime.Goexit()
	}
	printWarning()
	return uintptr(minusOne), 0, EACCES
}

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	if f := syscall("Syscall6"); f != nil {
		ch := chans.Get().(chan sysret)
		f.Invoke(func(r *js.Object) {
			ch <- sysret{uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())}
		}, trap, a1, a2, a3, a4, a5, a6)
		result := <-ch
		chans.Put(ch)
		return result.r1, result.r2, result.err
	}
	if trap != 202 { // kern.osrelease on OS X, happens in init of "os" package
		printWarning()
	}
	return uintptr(minusOne), 0, EACCES
}

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	return Syscall(trap, a1, a2, a3)
}

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	return Syscall6(trap, a1, a2, a3, a4, a5, a6)
}

func BytePtrFromString(s string) (*byte, error) {
	array := js.Global.Get("Uint8Array").New(len(s) + 1)
	for i, b := range []byte(s) {
		if b == 0 {
			return nil, EINVAL
		}
		array.SetIndex(i, b)
	}
	array.SetIndex(len(s), 0)
	return (*byte)(unsafe.Pointer(array.Unsafe())), nil
}

func forkExec(argv0 string, argv []string, attr *ProcAttr) (pid int, err error) {
	var err1 Errno

	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	// Convert args to C form.
	argv0p, err := BytePtrFromString(argv0)
	if err != nil {
		return 0, err
	}
	argvp, err := SlicePtrFromStrings(argv)
	if err != nil {
		return 0, err
	}
	envvp, err := SlicePtrFromStrings(attr.Env)
	if err != nil {
		return 0, err
	}

	if (runtime.GOOS == "freebsd" || runtime.GOOS == "dragonfly") && len(argv[0]) > len(argv0) {
		argvp[0] = argv0p
	}

	var chroot *byte
	chroot, err = BytePtrFromString(sys.Chroot)
	if err != nil {
		return 0, err
	}
	var dir *byte
	dir, err = BytePtrFromString(attr.Dir)
	if err != nil {
		return 0, err
	}

	// Acquire the fork lock so that no other threads
	// create new fds that are not yet close-on-exec
	// before we fork.
	ForkLock.Lock()

	// Kick off child.
	pid, err1 = forkAndExecInChild(argv0p, argvp, envvp, chroot, dir, attr, sys, 0)
	if err1 != 0 {
		err = Errno(err1)
		goto error
	}
	ForkLock.Unlock()

	// Read got EOF, so pipe closed on exec, so exec succeeded.
	return pid, nil

error:
	ForkLock.Unlock()
	return 0, err
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	// little to big
	sa.raw.Port = uint16(byte(sa.Port>>8)) + (uint16(sa.Port&0xff) << 8)
	//p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	//p[0] = byte(sa.Port >> 8)
	//p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	//up := unsafe.Pointer(&sa.raw)
	//return up, SizeofSockaddrInet4, nil
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
}

func forkAndExecInChild(argv0 *byte, argv, envv []*byte, chroot, dir *byte, attr *ProcAttr, sys *SysProcAttr, pipe int) (pid int, err Errno) {
	const SYS_SPAWN = 326

	sysPid, _, sysErr := RawSyscall6(
		SYS_SPAWN,
		uintptr(unsafe.Pointer(dir)),
		uintptr(unsafe.Pointer(argv0)),
		uintptr(unsafe.Pointer(&argv[0])),
		uintptr(unsafe.Pointer(&envv[0])),
		uintptr(unsafe.Pointer(&attr.Files[0])),
		uintptr(0),
	)

	return int(sysPid), Errno(sysErr)
}
