// +build js,!windows

package syscall

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
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

	// 	// Declare all variables at top in case any
	// 	// declarations require heap allocation (e.g., err1).
	// 	var (
	// 		r1     uintptr
	// 		err1   Errno
	// 		err2   Errno
	// 		nextfd int
	// 		i      int
	// 		p      [2]int
	// 	)

	// 	// Record parent PID so child can test if it has died.
	// 	ppid, _, _ := RawSyscall(SYS_GETPID, 0, 0, 0)

	// 	// Guard against side effects of shuffling fds below.
	// 	// Make sure that nextfd is beyond any currently open files so
	// 	// that we can't run the risk of overwriting any of them.
	// 	fd := make([]int, len(attr.Files))
	// 	nextfd = len(attr.Files)
	// 	for i, ufd := range attr.Files {
	// 		if nextfd < int(ufd) {
	// 			nextfd = int(ufd)
	// 		}
	// 		fd[i] = int(ufd)
	// 	}
	// 	nextfd++

	// 	// Allocate another pipe for parent to child communication for
	// 	// synchronizing writing of User ID/Group ID mappings.
	// 	if sys.UidMappings != nil || sys.GidMappings != nil {
	// 		if err := forkExecPipe(p[:]); err != nil {
	// 			return 0, err.(Errno)
	// 		}
	// 	}

	// 	// About to call fork.
	// 	// No more allocation or calls of non-assembly functions.
	// 	runtime_BeforeFork()
	// 	r1, _, err1 = RawSyscall6(SYS_CLONE, uintptr(SIGCHLD)|sys.Cloneflags, 0, 0, 0, 0, 0)
	// 	if err1 != 0 {
	// 		runtime_AfterFork()
	// 		return 0, err1
	// 	}

	// 	if r1 != 0 {
	// 		// parent; return PID
	// 		runtime_AfterFork()
	// 		pid = int(r1)

	// 		if sys.UidMappings != nil || sys.GidMappings != nil {
	// 			Close(p[0])
	// 			err := writeUidGidMappings(pid, sys)
	// 			if err != nil {
	// 				err2 = err.(Errno)
	// 			}
	// 			RawSyscall(SYS_WRITE, uintptr(p[1]), uintptr(unsafe.Pointer(&err2)), unsafe.Sizeof(err2))
	// 			Close(p[1])
	// 		}

	// 		return pid, 0
	// 	}

	// 	// Fork succeeded, now in child.

	// 	// Wait for User ID/Group ID mappings to be written.
	// 	if sys.UidMappings != nil || sys.GidMappings != nil {
	// 		if _, _, err1 = RawSyscall(SYS_CLOSE, uintptr(p[1]), 0, 0); err1 != 0 {
	// 			goto childerror
	// 		}
	// 		r1, _, err1 = RawSyscall(SYS_READ, uintptr(p[0]), uintptr(unsafe.Pointer(&err2)), unsafe.Sizeof(err2))
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 		if r1 != unsafe.Sizeof(err2) {
	// 			err1 = EINVAL
	// 			goto childerror
	// 		}
	// 		if err2 != 0 {
	// 			err1 = err2
	// 			goto childerror
	// 		}
	// 	}

	// 	// Enable tracing if requested.
	// 	if sys.Ptrace {
	// 		_, _, err1 = RawSyscall(SYS_PTRACE, uintptr(PTRACE_TRACEME), 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Session ID
	// 	if sys.Setsid {
	// 		_, _, err1 = RawSyscall(SYS_SETSID, 0, 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Set process group
	// 	if sys.Setpgid || sys.Foreground {
	// 		// Place child in process group.
	// 		_, _, err1 = RawSyscall(SYS_SETPGID, 0, uintptr(sys.Pgid), 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	if sys.Foreground {
	// 		pgrp := int32(sys.Pgid)
	// 		if pgrp == 0 {
	// 			r1, _, err1 = RawSyscall(SYS_GETPID, 0, 0, 0)
	// 			if err1 != 0 {
	// 				goto childerror
	// 			}

	// 			pgrp = int32(r1)
	// 		}

	// 		// Place process group in foreground.
	// 		_, _, err1 = RawSyscall(SYS_IOCTL, uintptr(sys.Ctty), uintptr(TIOCSPGRP), uintptr(unsafe.Pointer(&pgrp)))
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Chroot
	// 	if chroot != nil {
	// 		_, _, err1 = RawSyscall(SYS_CHROOT, uintptr(unsafe.Pointer(chroot)), 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// User and groups
	// 	if cred := sys.Credential; cred != nil {
	// 		ngroups := uintptr(len(cred.Groups))
	// 		if ngroups > 0 {
	// 			groups := unsafe.Pointer(&cred.Groups[0])
	// 			_, _, err1 = RawSyscall(SYS_SETGROUPS, ngroups, uintptr(groups), 0)
	// 			if err1 != 0 {
	// 				goto childerror
	// 			}
	// 		}
	// 		_, _, err1 = RawSyscall(SYS_SETGID, uintptr(cred.Gid), 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 		_, _, err1 = RawSyscall(SYS_SETUID, uintptr(cred.Uid), 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Chdir
	// 	if dir != nil {
	// 		_, _, err1 = RawSyscall(SYS_CHDIR, uintptr(unsafe.Pointer(dir)), 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Parent death signal
	// 	if sys.Pdeathsig != 0 {
	// 		_, _, err1 = RawSyscall6(SYS_PRCTL, PR_SET_PDEATHSIG, uintptr(sys.Pdeathsig), 0, 0, 0, 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}

	// 		// Signal self if parent is already dead. This might cause a
	// 		// duplicate signal in rare cases, but it won't matter when
	// 		// using SIGKILL.
	// 		r1, _, _ = RawSyscall(SYS_GETPPID, 0, 0, 0)
	// 		if r1 != ppid {
	// 			pid, _, _ := RawSyscall(SYS_GETPID, 0, 0, 0)
	// 			_, _, err1 := RawSyscall(SYS_KILL, pid, uintptr(sys.Pdeathsig), 0)
	// 			if err1 != 0 {
	// 				goto childerror
	// 			}
	// 		}
	// 	}

	// 	// Pass 1: look for fd[i] < i and move those up above len(fd)
	// 	// so that pass 2 won't stomp on an fd it needs later.
	// 	if pipe < nextfd {
	// 		_, _, err1 = RawSyscall(_SYS_dup, uintptr(pipe), uintptr(nextfd), 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 		RawSyscall(SYS_FCNTL, uintptr(nextfd), F_SETFD, FD_CLOEXEC)
	// 		pipe = nextfd
	// 		nextfd++
	// 	}
	// 	for i = 0; i < len(fd); i++ {
	// 		if fd[i] >= 0 && fd[i] < int(i) {
	// 			_, _, err1 = RawSyscall(_SYS_dup, uintptr(fd[i]), uintptr(nextfd), 0)
	// 			if err1 != 0 {
	// 				goto childerror
	// 			}
	// 			RawSyscall(SYS_FCNTL, uintptr(nextfd), F_SETFD, FD_CLOEXEC)
	// 			fd[i] = nextfd
	// 			nextfd++
	// 			if nextfd == pipe { // don't stomp on pipe
	// 				nextfd++
	// 			}
	// 		}
	// 	}

	// 	// Pass 2: dup fd[i] down onto i.
	// 	for i = 0; i < len(fd); i++ {
	// 		if fd[i] == -1 {
	// 			RawSyscall(SYS_CLOSE, uintptr(i), 0, 0)
	// 			continue
	// 		}
	// 		if fd[i] == int(i) {
	// 			// dup2(i, i) won't clear close-on-exec flag on Linux,
	// 			// probably not elsewhere either.
	// 			_, _, err1 = RawSyscall(SYS_FCNTL, uintptr(fd[i]), F_SETFD, 0)
	// 			if err1 != 0 {
	// 				goto childerror
	// 			}
	// 			continue
	// 		}
	// 		// The new fd is created NOT close-on-exec,
	// 		// which is exactly what we want.
	// 		_, _, err1 = RawSyscall(_SYS_dup, uintptr(fd[i]), uintptr(i), 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// By convention, we don't close-on-exec the fds we are
	// 	// started with, so if len(fd) < 3, close 0, 1, 2 as needed.
	// 	// Programs that know they inherit fds >= 3 will need
	// 	// to set them close-on-exec.
	// 	for i = len(fd); i < 3; i++ {
	// 		RawSyscall(SYS_CLOSE, uintptr(i), 0, 0)
	// 	}

	// 	// Detach fd 0 from tty
	// 	if sys.Noctty {
	// 		_, _, err1 = RawSyscall(SYS_IOCTL, 0, uintptr(TIOCNOTTY), 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Set the controlling TTY to Ctty
	// 	if sys.Setctty {
	// 		_, _, err1 = RawSyscall(SYS_IOCTL, uintptr(sys.Ctty), uintptr(TIOCSCTTY), 0)
	// 		if err1 != 0 {
	// 			goto childerror
	// 		}
	// 	}

	// 	// Time to exec.
	// 	_, _, err1 = RawSyscall(SYS_EXECVE,
	// 		uintptr(unsafe.Pointer(argv0)),
	// 		uintptr(unsafe.Pointer(&argv[0])),
	// 		uintptr(unsafe.Pointer(&envv[0])))

	// childerror:
	// 	// send error code on pipe
	// 	RawSyscall(SYS_WRITE, uintptr(pipe), uintptr(unsafe.Pointer(&err1)), unsafe.Sizeof(err1))
	// 	for {
	// 		RawSyscall(SYS_EXIT, 253, 0, 0)
	// 	}
}
