// AUTO GENERATED - DO NOT EDIT

package prelude

const syscalls = `
(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
'use strict';
(function (AF) {
    AF[AF["UNSPEC"] = 0] = "UNSPEC";
    AF[AF["LOCAL"] = 1] = "LOCAL";
    AF[AF["UNIX"] = 1] = "UNIX";
    AF[AF["FILE"] = 1] = "FILE";
    AF[AF["INET"] = 2] = "INET";
    AF[AF["INET6"] = 10] = "INET6";
})(exports.AF || (exports.AF = {}));
var AF = exports.AF;
;
(function (SOCK) {
    SOCK[SOCK["STREAM"] = 1] = "STREAM";
    SOCK[SOCK["DGRAM"] = 2] = "DGRAM";
})(exports.SOCK || (exports.SOCK = {}));
var SOCK = exports.SOCK;
(function (ErrorCode) {
    ErrorCode[ErrorCode["EPERM"] = 0] = "EPERM";
    ErrorCode[ErrorCode["ENOENT"] = 1] = "ENOENT";
    ErrorCode[ErrorCode["EIO"] = 2] = "EIO";
    ErrorCode[ErrorCode["EBADF"] = 3] = "EBADF";
    ErrorCode[ErrorCode["EACCES"] = 4] = "EACCES";
    ErrorCode[ErrorCode["EBUSY"] = 5] = "EBUSY";
    ErrorCode[ErrorCode["EEXIST"] = 6] = "EEXIST";
    ErrorCode[ErrorCode["ENOTDIR"] = 7] = "ENOTDIR";
    ErrorCode[ErrorCode["EISDIR"] = 8] = "EISDIR";
    ErrorCode[ErrorCode["EINVAL"] = 9] = "EINVAL";
    ErrorCode[ErrorCode["EFBIG"] = 10] = "EFBIG";
    ErrorCode[ErrorCode["ENOSPC"] = 11] = "ENOSPC";
    ErrorCode[ErrorCode["EROFS"] = 12] = "EROFS";
    ErrorCode[ErrorCode["ENOTEMPTY"] = 13] = "ENOTEMPTY";
    ErrorCode[ErrorCode["ENOTSUP"] = 14] = "ENOTSUP";
})(exports.ErrorCode || (exports.ErrorCode = {}));
var ErrorCode = exports.ErrorCode;
var fsErrors = {
    EPERM: 'Operation not permitted.',
    ENOENT: 'No such file or directory.',
    EIO: 'Input/output error.',
    EBADF: 'Bad file descriptor.',
    EACCES: 'Permission denied.',
    EBUSY: 'Resource busy or locked.',
    EEXIST: 'File exists.',
    ENOTDIR: 'File is not a directory.',
    EISDIR: 'File is a directory.',
    EINVAL: 'Invalid argument.',
    EFBIG: 'File is too big.',
    ENOSPC: 'No space left on disk.',
    EROFS: 'Cannot modify a read-only file system.',
    ENOTEMPTY: 'Directory is not empty.',
    ENOTSUP: 'Operation is not supported.',
};
var ApiError = (function () {
    function ApiError(type, message) {
        this.type = type;
        this.code = ErrorCode[type];
        if (message != null) {
            this.message = message;
        }
        else {
            this.message = fsErrors[type];
        }
    }
    ApiError.prototype.toString = function () {
        return this.code + ": " + fsErrors[this.code] + " " + this.message;
    };
    return ApiError;
})();
exports.ApiError = ApiError;
function convertApiErrors(e) {
    if (!e)
        return e;
    if (!e.hasOwnProperty('type') || !e.hasOwnProperty('message') || !e.hasOwnProperty('code'))
        return e;
    return new ApiError(e.type, e.message);
}
var SyscallResponse = (function () {
    function SyscallResponse(id, name, args) {
        this.id = id;
        this.name = name;
        this.args = args;
    }
    SyscallResponse.From = function (ev) {
        if (!ev.data)
            return;
        for (var i = 0; i < SyscallResponse.requiredOnData.length; i++) {
            if (!ev.data.hasOwnProperty(SyscallResponse.requiredOnData[i]))
                return;
        }
        var args = ev.data.args.map(convertApiErrors);
        return new SyscallResponse(ev.data.id, ev.data.name, args);
    };
    SyscallResponse.requiredOnData = ['id', 'name', 'args'];
    return SyscallResponse;
})();
exports.SyscallResponse = SyscallResponse;
var USyscalls = (function () {
    function USyscalls(port) {
        this.msgIdSeq = 1;
        this.outstanding = {};
        this.signalHandlers = {};
        this.port = port;
        this.port.onmessage = this.resultHandler.bind(this);
    }
    USyscalls.prototype.exit = function (code) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = function () {
            var args = [];
            for (var _i = 0; _i < arguments.length; _i++) {
                args[_i - 0] = arguments[_i];
            }
            console.log('received callback for exit(), should clean up');
        };
        this.post(msgId, 'exit', code);
    };
    USyscalls.prototype.kill = function (pid, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'kill', pid);
    };
    USyscalls.prototype.socket = function (domain, type, protocol, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'socket', domain, type, protocol);
    };
    USyscalls.prototype.getsockname = function (fd, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'getsockname', fd);
    };
    USyscalls.prototype.bind = function (fd, sockInfo, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'bind', fd, sockInfo);
    };
    USyscalls.prototype.listen = function (fd, backlog, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'listen', fd, backlog);
    };
    USyscalls.prototype.accept = function (fd, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'accept', fd);
    };
    USyscalls.prototype.connect = function (fd, addr, port, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'connect', fd, addr, port);
    };
    USyscalls.prototype.getcwd = function (cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'getcwd');
    };
    USyscalls.prototype.getpid = function (cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'getpid');
    };
    USyscalls.prototype.spawn = function (cwd, name, args, env, files, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'spawn', cwd, name, args, env, files);
    };
    USyscalls.prototype.pipe2 = function (flags, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'pipe2', flags);
    };
    USyscalls.prototype.getpriority = function (which, who, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'getpriority', which, who);
    };
    USyscalls.prototype.setpriority = function (which, who, prio, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'setpriority', which, who, prio);
    };
    USyscalls.prototype.open = function (path, flags, mode, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'open', path, flags, mode);
    };
    USyscalls.prototype.unlink = function (path, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'unlink', path);
    };
    USyscalls.prototype.utimes = function (path, atime, mtime, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'utimes', path, atime, mtime);
    };
    USyscalls.prototype.futimes = function (fd, atime, mtime, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'futimes', fd, atime, mtime);
    };
    USyscalls.prototype.rmdir = function (path, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'rmdir', path);
    };
    USyscalls.prototype.mkdir = function (path, mode, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'mkdir', path);
    };
    USyscalls.prototype.close = function (fd, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'close', fd);
    };
    USyscalls.prototype.pwrite = function (fd, buf, pos, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'pwrite', fd, buf, pos);
    };
    USyscalls.prototype.readdir = function (path, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'readdir', path);
    };
    USyscalls.prototype.fstat = function (fd, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'fstat', fd);
    };
    USyscalls.prototype.lstat = function (path, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'lstat', path);
    };
    USyscalls.prototype.stat = function (path, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'stat', path);
    };
    USyscalls.prototype.ioctl = function (fd, request, length, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'ioctl', fd, request, length);
    };
    USyscalls.prototype.readlink = function (path, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'readlink', path);
    };
    USyscalls.prototype.getdents = function (fd, length, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'getdents', fd, length);
    };
    USyscalls.prototype.pread = function (fd, length, offset, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'pread', fd, length, offset);
    };
    USyscalls.prototype.addEventListener = function (type, handler) {
        if (!handler)
            return;
        if (this.signalHandlers[type])
            this.signalHandlers[type].push(handler);
        else
            this.signalHandlers[type] = [handler];
    };
    USyscalls.prototype.resultHandler = function (ev) {
        var response = SyscallResponse.From(ev);
        if (!response) {
            console.log('bad usyscall message, dropping');
            console.log(ev);
            return;
        }
        if (response.name) {
            var handlers = this.signalHandlers[response.name];
            if (handlers) {
                for (var i = 0; i < handlers.length; i++)
                    handlers[i](response);
            }
            else {
                console.log('unhandled signal ' + response.name);
            }
            return;
        }
        this.complete(response.id, response.args);
    };
    USyscalls.prototype.complete = function (id, args) {
        var cb = this.outstanding[id];
        delete this.outstanding[id];
        if (cb) {
            cb.apply(undefined, args);
        }
        else {
            console.log('unknown callback for msg ' + id + ' - ' + args);
        }
    };
    USyscalls.prototype.nextMsgId = function () {
        return ++this.msgIdSeq;
    };
    USyscalls.prototype.post = function (msgId, name) {
        var args = [];
        for (var _i = 2; _i < arguments.length; _i++) {
            args[_i - 2] = arguments[_i];
        }
        this.port.postMessage({
            id: msgId,
            name: name,
            args: args,
        });
    };
    return USyscalls;
})();
exports.USyscalls = USyscalls;
function getGlobal() {
    if (typeof window !== "undefined") {
        return window;
    }
    else if (typeof self !== "undefined") {
        return self;
    }
    else if (typeof global !== "undefined") {
        return global;
    }
    else {
        return this;
    }
}
exports.getGlobal = getGlobal;
exports.syscall = new USyscalls(getGlobal());

},{}],2:[function(require,module,exports){
'use strict';
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var syscall_1 = require('../browser-node/syscall');
var table_1 = require('./table');
function Syscall(cb, trap) {
    table_1.syscallTbl[trap].apply(this, arguments);
}
exports.Syscall = Syscall;
;
exports.Syscall6 = Syscall;
exports.internal = syscall_1.syscall;
var OnceEmitter = (function () {
    function OnceEmitter() {
        this.listeners = {};
    }
    OnceEmitter.prototype.once = function (event, cb) {
        var cbs = this.listeners[event];
        if (!cbs)
            cbs = [cb];
        else
            cbs.push(cb);
        this.listeners[event] = cbs;
    };
    OnceEmitter.prototype.emit = function (event) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        var cbs = this.listeners[event];
        this.listeners[event] = [];
        if (!cbs)
            return;
        for (var i = 0; i < cbs.length; i++) {
            cbs[i].apply(null, args);
        }
    };
    return OnceEmitter;
})();
var Process = (function (_super) {
    __extends(Process, _super);
    function Process(argv, environ) {
        _super.call(this);
        this.argv = argv;
        this.env = environ;
    }
    Process.prototype.exit = function (code) {
        if (code === void 0) { code = 0; }
        syscall_1.syscall.exit(code);
    };
    return Process;
})(OnceEmitter);
var process = new Process(null, null);
syscall_1.syscall.addEventListener('init', init.bind(this));
function init(data) {
    'use strict';
    var args = data.args[0];
    var environ = data.args[1];
    args = [args[0]].concat(args);
    process.argv = args;
    process.env = environ;
    setTimeout(function () { process.emit('ready'); }, 0);
}
if (typeof window !== "undefined") {
    window.$syscall = exports;
    window.process = process;
}
else if (typeof self !== "undefined") {
    self.$syscall = exports;
    self.process = process;
}
else if (typeof global !== "undefined") {
    global.$syscall = exports;
}
else {
    this.$syscall = exports;
    this.process = process;
}

},{"../browser-node/syscall":1,"./table":3}],3:[function(require,module,exports){
'use strict';
var syscall_1 = require('../browser-node/syscall');
var ENOSYS = 38;
var AT_FDCWD = -0x64;
function sys_ni_syscall(cb, trap) {
    console.log('ni syscall ' + trap);
    debugger;
    setTimeout(cb, 0, [-1, 0, -ENOSYS]);
}
function sys_getpid(cb, trap) {
    var done = function (err, pid) {
        cb([pid, 0, 0]);
    };
    syscall_1.syscall.getpid.apply(syscall_1.syscall, [done]);
}
function sys_getcwd(cb, trap, arg0, arg1, arg2) {
    var $getcwdArray = arg0;
    var $getcwdLen = arg1;
    var done = function (p) {
        for (var i = 0; i < p.length; i++)
            $getcwdArray[i] = p.charCodeAt(i);
        var nullPos = p.length;
        if (nullPos >= $getcwdArray.byteLength)
            nullPos = $getcwdArray.byteLength;
        $getcwdArray[nullPos] = 0;
        cb([p.length + 1, 0, 0]);
    };
    syscall_1.syscall.getcwd.apply(syscall_1.syscall, [done]);
}
function sys_ioctl(cb, trap, arg0, arg1, arg2) {
    var $fd = arg0;
    var $request = arg1;
    var $argp = arg2;
    var done = function (err, buf) {
        if (!err && $argp.byteLength !== undefined)
            $argp.set(buf);
        cb([err ? err : buf.byteLength, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.ioctl.apply(syscall_1.syscall, [$fd, $request, $argp.byteLength, done]);
}
function sys_getdents64(cb, trap, arg0, arg1, arg2) {
    var $fd = arg0;
    var $buf = arg1;
    var $len = arg2;
    var done = function (err, buf) {
        if (!err)
            $buf.set(buf);
        cb([err ? -1 : buf.byteLength, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.getdents.apply(syscall_1.syscall, [$fd, $len, done]);
}
function sys_read(cb, trap, arg0, arg1, arg2) {
    var $readArray = arg1;
    var $readLen = arg2;
    var done = function (err, dataLen, data) {
        if (!err) {
            for (var i = 0; i < dataLen; i++)
                $readArray[i] = data[i];
        }
        cb([dataLen, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.pread.apply(syscall_1.syscall, [arg0, arg2, -1, done]);
}
function sys_write(cb, trap, arg0, arg1, arg2) {
    var done = function (err, len) {
        cb([len, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.pwrite.apply(syscall_1.syscall, [arg0, new Uint8Array(arg1, 0, arg2), 0, done]);
}
function sys_stat(cb, trap, arg0, arg1) {
    var $fstatArray = arg1;
    var done = function (err, stats) {
        if (!err)
            $fstatArray.set(stats);
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    var len = arg0.length;
    if (len && arg0[arg0.length - 1] === 0)
        len--;
    syscall_1.syscall.stat.apply(syscall_1.syscall, [arg0.slice(0, len), done]);
}
function sys_lstat(cb, trap, arg0, arg1) {
    var $fstatArray = arg1;
    var done = function (err, stats) {
        if (!err)
            $fstatArray.set(stats);
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    var len = arg0.length;
    if (len && arg0[arg0.length - 1] === 0)
        len--;
    syscall_1.syscall.lstat.apply(syscall_1.syscall, [arg0.slice(0, len), done]);
}
function sys_fstat(cb, trap, arg0, arg1) {
    var $fstatArray = arg1;
    var done = function (err, stats) {
        if (!err)
            $fstatArray.set(stats);
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.fstat.apply(syscall_1.syscall, [arg0, done]);
}
function sys_readlinkat(cb, trap, arg0, arg1, arg2, arg3) {
    var $fd = arg0 | 0;
    var $path = arg1;
    var $buf = arg2;
    var $buflen = arg3;
    if ((arg0 | 0) !== AT_FDCWD) {
        debugger;
        setTimeout(cb, 0, [-1, 0, -1]);
        return;
    }
    var done = function (err, linkString) {
        if (!err)
            $buf.set(linkString);
        cb([err ? -1 : linkString.length, 0, err ? -1 : 0]);
    };
    var len = $path.length;
    if (len && arg1[$path.length - 1] === 0)
        len--;
    syscall_1.syscall.readlink.apply(syscall_1.syscall, [$path.slice(0, len), done]);
}
function sys_openat(cb, trap, arg0, arg1, arg2, arg3) {
    if ((arg0 | 0) !== AT_FDCWD) {
        debugger;
        setTimeout(cb, 0, [-1, 0, -1]);
        return;
    }
    var done = function (err, fd) {
        if (err)
            console.log('error: ' + err);
        cb([fd, 0, err ? -1 : 0]);
    };
    var len = arg1.length;
    if (len && arg1[arg1.length - 1] === 0)
        len--;
    syscall_1.syscall.open.apply(syscall_1.syscall, [arg1.slice(0, len), arg2, arg3, done]);
}
function sys_close(cb, trap, arg0) {
    var done = function (err) {
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.close.apply(syscall_1.syscall, [arg0, done]);
}
function sys_exit_group(cb, trap, arg0) {
    syscall_1.syscall.exit(arg0);
}
function sys_socket(cb, trap, arg0, arg1, arg2) {
    var done = function (err, fd) {
        cb([err ? -1 : fd, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.socket(arg0, arg1, arg2, done);
}
function sys_bind(cb, trap, arg0, arg1, arg2) {
    console.log('FIXME: unmarshal');
    var done = function (err) {
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.bind(arg0, arg1.slice(0, arg2), done);
}
function sys_listen(cb, trap, arg0, arg1) {
    var done = function (err) {
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.listen(arg0, arg1, done);
}
function sys_getsockname(cb, trap, arg0, arg1, arg2) {
    var done = function (err, buf) {
        if (!err) {
            arg1.set(buf);
            arg2.$set(buf.byteLength);
        }
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.getsockname(arg0, done);
}
function sys_accept4(cb, trap, arg0, arg1, arg2) {
    var $acceptArray = arg1;
    var $acceptLen = arg2;
    var done = function (err, fd, sockInfo) {
        $acceptArray.set(sockInfo);
        $acceptLen.$set(sockInfo.length);
        cb([err ? -1 : fd, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.accept(arg0, done);
}
function sys_setsockopt(cb, trap) {
    console.log('FIXME: implement setsockopt');
    setTimeout(cb, 0, [0, 0, 0]);
}
exports.syscallTbl = [
    sys_read,
    sys_write,
    sys_ni_syscall,
    sys_close,
    sys_stat,
    sys_fstat,
    sys_lstat,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ioctl,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_getpid,
    sys_ni_syscall,
    sys_socket,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_bind,
    sys_listen,
    sys_getsockname,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_setsockopt,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_getcwd,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_getdents64,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_exit_group,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_openat,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_readlinkat,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_accept4,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
    sys_ni_syscall,
];

},{"../browser-node/syscall":1}]},{},[2]);
`
