// AUTO GENERATED - DO NOT EDIT

package prelude

const syscalls = `
(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
'use strict';
exports.kMaxLength = 0x3fffffff;
function blitBuffer(src, dst, offset, length) {
    var i;
    for (i = 0; i < length; i++) {
        if ((i + offset >= dst.length) || (i >= src.length))
            break;
        dst[i + offset] = src[i];
    }
    return i;
}
function utf8Slice(buf, start, end) {
    end = Math.min(buf.length, end);
    var res = [];
    var i = start;
    while (i < end) {
        var firstByte = buf[i];
        var codePoint = null;
        var bytesPerSequence = (firstByte > 0xEF) ? 4
            : (firstByte > 0xDF) ? 3
                : (firstByte > 0xBF) ? 2
                    : 1;
        if (i + bytesPerSequence <= end) {
            var secondByte = void 0, thirdByte = void 0, fourthByte = void 0, tempCodePoint = void 0;
            switch (bytesPerSequence) {
                case 1:
                    if (firstByte < 0x80) {
                        codePoint = firstByte;
                    }
                    break;
                case 2:
                    secondByte = buf[i + 1];
                    if ((secondByte & 0xC0) === 0x80) {
                        tempCodePoint = (firstByte & 0x1F) << 0x6 | (secondByte & 0x3F);
                        if (tempCodePoint > 0x7F) {
                            codePoint = tempCodePoint;
                        }
                    }
                    break;
                case 3:
                    secondByte = buf[i + 1];
                    thirdByte = buf[i + 2];
                    if ((secondByte & 0xC0) === 0x80 && (thirdByte & 0xC0) === 0x80) {
                        tempCodePoint = (firstByte & 0xF) << 0xC | (secondByte & 0x3F) << 0x6 | (thirdByte & 0x3F);
                        if (tempCodePoint > 0x7FF && (tempCodePoint < 0xD800 || tempCodePoint > 0xDFFF)) {
                            codePoint = tempCodePoint;
                        }
                    }
                    break;
                case 4:
                    secondByte = buf[i + 1];
                    thirdByte = buf[i + 2];
                    fourthByte = buf[i + 3];
                    if ((secondByte & 0xC0) === 0x80 && (thirdByte & 0xC0) === 0x80 && (fourthByte & 0xC0) === 0x80) {
                        tempCodePoint = (firstByte & 0xF) << 0x12 | (secondByte & 0x3F) << 0xC | (thirdByte & 0x3F) << 0x6 | (fourthByte & 0x3F);
                        if (tempCodePoint > 0xFFFF && tempCodePoint < 0x110000) {
                            codePoint = tempCodePoint;
                        }
                    }
            }
        }
        if (codePoint === null) {
            codePoint = 0xFFFD;
            bytesPerSequence = 1;
        }
        else if (codePoint > 0xFFFF) {
            codePoint -= 0x10000;
            res.push(codePoint >>> 10 & 0x3FF | 0xD800);
            codePoint = 0xDC00 | codePoint & 0x3FF;
        }
        res.push(codePoint);
        i += bytesPerSequence;
    }
    return decodeCodePointsArray(res);
}
exports.utf8Slice = utf8Slice;
var MAX_ARGUMENTS_LENGTH = 0x1000;
function decodeCodePointsArray(codePoints) {
    var len = codePoints.length;
    if (len <= MAX_ARGUMENTS_LENGTH) {
        return String.fromCharCode.apply(String, codePoints);
    }
    var res = '';
    var i = 0;
    while (i < len) {
        res += String.fromCharCode.apply(String, codePoints.slice(i, i += MAX_ARGUMENTS_LENGTH));
    }
    return res;
}
function utf8ToBytes(string, units) {
    units = units || Infinity;
    var codePoint;
    var length = string.length;
    var leadSurrogate = null;
    var bytes = [];
    for (var i = 0; i < length; i++) {
        codePoint = string.charCodeAt(i);
        if (codePoint > 0xD7FF && codePoint < 0xE000) {
            if (!leadSurrogate) {
                if (codePoint > 0xDBFF) {
                    if ((units -= 3) > -1)
                        bytes.push(0xEF, 0xBF, 0xBD);
                    continue;
                }
                else if (i + 1 === length) {
                    if ((units -= 3) > -1)
                        bytes.push(0xEF, 0xBF, 0xBD);
                    continue;
                }
                leadSurrogate = codePoint;
                continue;
            }
            if (codePoint < 0xDC00) {
                if ((units -= 3) > -1)
                    bytes.push(0xEF, 0xBF, 0xBD);
                leadSurrogate = codePoint;
                continue;
            }
            codePoint = leadSurrogate - 0xD800 << 10 | codePoint - 0xDC00 | 0x10000;
        }
        else if (leadSurrogate) {
            if ((units -= 3) > -1)
                bytes.push(0xEF, 0xBF, 0xBD);
        }
        leadSurrogate = null;
        if (codePoint < 0x80) {
            if ((units -= 1) < 0)
                break;
            bytes.push(codePoint);
        }
        else if (codePoint < 0x800) {
            if ((units -= 2) < 0)
                break;
            bytes.push(codePoint >> 0x6 | 0xC0, codePoint & 0x3F | 0x80);
        }
        else if (codePoint < 0x10000) {
            if ((units -= 3) < 0)
                break;
            bytes.push(codePoint >> 0xC | 0xE0, codePoint >> 0x6 & 0x3F | 0x80, codePoint & 0x3F | 0x80);
        }
        else if (codePoint < 0x110000) {
            if ((units -= 4) < 0)
                break;
            bytes.push(codePoint >> 0x12 | 0xF0, codePoint >> 0xC & 0x3F | 0x80, codePoint >> 0x6 & 0x3F | 0x80, codePoint & 0x3F | 0x80);
        }
        else {
            throw new Error('Invalid code point');
        }
    }
    return bytes;
}
function asciiSlice(buf, start, end) {
    var ret = '';
    end = Math.min(buf.length, end);
    for (var i = start; i < end; i++) {
        ret += String.fromCharCode(buf[i] & 0x7F);
    }
    return ret;
}
function setupBufferJS(prototype, bindingObj) {
    bindingObj.flags = [0];
    prototype.__proto__ = Uint8Array.prototype;
    prototype.utf8Write = function (str, offset, length) {
        return blitBuffer(utf8ToBytes(str, this.length - offset), this, offset, length);
    };
    prototype.utf8Slice = function (start, end) {
        return utf8Slice(this, start, end);
    };
    prototype.asciiSlice = function (start, end) {
        return asciiSlice(this, start, end);
    };
    prototype.copy = function copy(target, targetStart, start, end) {
        if (!start)
            start = 0;
        if (!end && end !== 0)
            end = this.length;
        if (targetStart >= target.length)
            targetStart = target.length;
        if (!targetStart)
            targetStart = 0;
        if (end > 0 && end < start)
            end = start;
        if (end === start)
            return 0;
        if (target.length === 0 || this.length === 0)
            return 0;
        if (targetStart < 0) {
            throw new RangeError('targetStart out of bounds');
        }
        if (start < 0 || start >= this.length)
            throw new RangeError('sourceStart out of bounds');
        if (end < 0)
            throw new RangeError('sourceEnd out of bounds');
        if (end > this.length)
            end = this.length;
        if (target.length - targetStart < end - start) {
            end = target.length - targetStart + start;
        }
        var len = end - start;
        var i;
        if (this === target && start < targetStart && targetStart < end) {
            for (i = len - 1; i >= 0; i--) {
                target[i + targetStart] = this[i + start];
            }
        }
        else {
            for (i = 0; i < len; i++) {
                target[i + targetStart] = this[i + start];
            }
        }
        return len;
    };
}
exports.setupBufferJS = setupBufferJS;
function createFromString(str, encoding) {
    console.log('TODO: createFromString');
}
exports.createFromString = createFromString;
function createFromArrayBuffer(obj) {
    console.log('TODO: createFromArrayBuffer');
}
exports.createFromArrayBuffer = createFromArrayBuffer;
function compare(a, b) {
    console.log('TODO: compare');
}
exports.compare = compare;
function byteLengthUtf8(str) {
    return utf8ToBytes(str).length;
}
exports.byteLengthUtf8 = byteLengthUtf8;
function indexOfString(buf, val, byteOffset) {
    console.log('TODO: indexOfString');
}
exports.indexOfString = indexOfString;
function indexOfBuffer(buf, val, byteOffset) {
    console.log('TODO: indexOfBuffer');
}
exports.indexOfBuffer = indexOfBuffer;
function indexOfNumber(buf, val, byteOffset) {
    console.log('TODO: indexOfNumber');
}
exports.indexOfNumber = indexOfNumber;
function fill(buf, val, start, end) {
    console.log('TODO: fill');
}
exports.fill = fill;
function readFloatLE(buf, offset) {
    console.log('TODO: readFloatLE');
}
exports.readFloatLE = readFloatLE;
function readFloatBE(buf, offset) {
    console.log('TODO: readFloatBE');
}
exports.readFloatBE = readFloatBE;
function readDoubleLE(buf, offset) {
    console.log('TODO: readDoubleLE');
}
exports.readDoubleLE = readDoubleLE;
function readDoubleBE(buf, offset) {
    console.log('TODO: readDoubleBE');
}
exports.readDoubleBE = readDoubleBE;
function writeFloatLE(buf, val, offset) {
    console.log('TODO: writeFloatLE');
}
exports.writeFloatLE = writeFloatLE;
function writeFloatBE(buf, val, offset) {
    console.log('TODO: writeFloatBE');
}
exports.writeFloatBE = writeFloatBE;
function writeDoubleLE(buf, val, offset) {
    console.log('TODO: writeDoubleLE');
}
exports.writeDoubleLE = writeDoubleLE;
function writeDoubleBE(buf, val, offset) {
    console.log('TODO: writeDoubleBE');
}
exports.writeDoubleBE = writeDoubleBE;

},{}],2:[function(require,module,exports){
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
        this.syscallPending = false;
        this.msgQueue = [];
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
    USyscalls.prototype.bind = function (fd, addr, port, cb) {
        var msgId = this.nextMsgId();
        this.outstanding[msgId] = cb;
        this.post(msgId, 'bind', fd, addr, port);
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
        this.syscallPending = false;
        if (this.msgQueue.length)
            setTimeout(this.doPost.bind(this), 0);
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
        this.msgQueue.push({
            id: msgId,
            name: name,
            args: args,
        });
        setTimeout(this.doPost.bind(this), 0);
    };
    USyscalls.prototype.doPost = function () {
        for (var msg = this.msgQueue.shift(); msg; msg = this.msgQueue.shift()) {
            this.port.postMessage(msg);
        }
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

},{}],3:[function(require,module,exports){
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

},{"../browser-node/syscall":2,"./table":4}],4:[function(require,module,exports){
'use strict';
var node_binary_marshal_1 = require('node-binary-marshal');
var syscall_1 = require('../browser-node/syscall');
var buffer_1 = require('../browser-node/binding/buffer');
var ENOSYS = 38;
var AT_FDCWD = -0x64;
function sys_ni_syscall(cb, trap) {
    console.log('ni syscall ' + trap);
    debugger;
    setTimeout(cb, 0, [-1, 0, -ENOSYS]);
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
    debugger;
    var done = function (err, buf) {
        if (!err)
            $argp.set(buf);
        cb([err ? -1 : buf.byteLength, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.ioctl.apply(syscall_1.syscall, [$fd, $request, $argp.byteLength]);
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
        var view = new DataView($fstatArray.buffer, $fstatArray.byteOffset);
        node_binary_marshal_1.Marshal(view, 0, stats, node_binary_marshal_1.fs.StatDef);
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    var len = arg0.length;
    if (len && arg0[arg0.length - 1] === 0)
        len--;
    var s = buffer_1.utf8Slice(arg0, 0, len);
    syscall_1.syscall.stat.apply(syscall_1.syscall, [s, done]);
}
function sys_lstat(cb, trap, arg0, arg1) {
    var $fstatArray = arg1;
    var done = function (err, stats) {
        var view = new DataView($fstatArray.buffer, $fstatArray.byteOffset);
        node_binary_marshal_1.Marshal(view, 0, stats, node_binary_marshal_1.fs.StatDef);
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    var len = arg0.length;
    if (len && arg0[arg0.length - 1] === 0)
        len--;
    var s = buffer_1.utf8Slice(arg0, 0, len);
    syscall_1.syscall.lstat.apply(syscall_1.syscall, [s, done]);
}
function sys_fstat(cb, trap, arg0, arg1) {
    var $fstatArray = arg1;
    var done = function (err, stats) {
        var view = new DataView($fstatArray.buffer, $fstatArray.byteOffset);
        node_binary_marshal_1.Marshal(view, 0, stats, node_binary_marshal_1.fs.StatDef);
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.fstat.apply(syscall_1.syscall, [arg0, done]);
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
    syscall_1.syscall.open.apply(syscall_1.syscall, [buffer_1.utf8Slice(arg1, 0, len), arg2, arg3, done]);
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
    syscall_1.syscall.socket.apply(syscall_1.syscall, [arg0, arg1, arg2, done]);
}
function sys_bind(cb, trap, arg0, arg1, arg2) {
    console.log('FIXME: unmarshal');
    var done = function (err) {
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.bind.apply(syscall_1.syscall, [arg0, '127.0.0.1', 8080, done]);
}
function sys_listen(cb, trap, arg0, arg1) {
    var done = function (err) {
        cb([err ? -1 : 0, 0, err ? -1 : 0]);
    };
    syscall_1.syscall.listen.apply(syscall_1.syscall, [arg0, arg1, done]);
}
function sys_getsockname(cb, trap, arg0, arg1, arg2) {
    console.log('TODO: getsockname');
    var view = new DataView(arg1.buffer, arg1.byteOffset);
    node_binary_marshal_1.Marshal(arg1, 0, { family: 2, port: 8080, addr: '127.0.0.1' }, node_binary_marshal_1.socket.SockAddrInDef);
    arg2.$set(node_binary_marshal_1.socket.SockAddrInDef.length);
    setTimeout(cb, 0, [0, 0, 0]);
}
function sys_accept4(cb, trap, arg0, arg1, arg2) {
    var $acceptArray = arg1;
    var $acceptLen = arg2;
    var args = [arg0, function (err, fd, remoteAddr, remotePort) {
            if (remoteAddr === 'localhost')
                remoteAddr = '127.0.0.1';
            var view = new DataView($acceptArray.buffer, $acceptArray.byteOffset);
            node_binary_marshal_1.Marshal(view, 0, { family: 2, port: remotePort, addr: remoteAddr }, node_binary_marshal_1.socket.SockAddrInDef);
            $acceptLen.$set(node_binary_marshal_1.socket.SockAddrInDef.length);
            cb([err ? -1 : fd, 0, err ? -1 : 0]);
        }];
    syscall_1.syscall.accept.apply(syscall_1.syscall, args);
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
    sys_ni_syscall,
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

},{"../browser-node/binding/buffer":1,"../browser-node/syscall":2,"node-binary-marshal":6}],5:[function(require,module,exports){
'use strict';
var marshal_1 = require('./marshal');
var utf8 = require('./utf8');
exports.TimespecDef = {
    fields: [
        { name: 'sec', type: 'int64' },
        { name: 'nsec', type: 'int64' },
    ],
    alignment: 'natural',
    length: 16,
};
exports.TimevalDef = exports.TimespecDef;
function nullMarshal(dst, off, src) {
    return [undefined, null];
}
;
function nullUnmarshal(src, off) {
    return [null, undefined, null];
}
;
function timespecMarshal(dst, off, src) {
    var timestamp = Date.parse(src);
    var secs = Math.floor(timestamp / 1000);
    var timespec = {
        sec: secs,
        nsec: (timestamp - secs * 1000) * 1e6,
    };
    return marshal_1.Marshal(dst, off, timespec, exports.TimespecDef);
}
;
function timespecUnmarshal(src, off) {
    var timespec = {};
    var _a = marshal_1.Unmarshal(timespec, src, off, exports.TimespecDef), len = _a[0], err = _a[1];
    var sec = timespec.sec;
    var nsec = timespec.nsec;
    var timestr = new Date(sec * 1e3 + nsec / 1e6).toISOString();
    return [timestr, len, err];
}
;
exports.StatDef = {
    fields: [
        { name: 'dev', type: 'uint64' },
        { name: 'ino', type: 'uint64' },
        { name: 'nlink', type: 'uint64' },
        { name: 'mode', type: 'uint32' },
        { name: 'uid', type: 'uint32' },
        { name: 'gid', type: 'uint32' },
        { name: 'X__pad0', type: 'int32', marshal: nullMarshal, unmarshal: nullUnmarshal, omit: true },
        { name: 'rdev', type: 'uint64' },
        { name: 'size', type: 'int64' },
        { name: 'blksize', type: 'int64' },
        { name: 'blocks', type: 'int64' },
        { name: 'atime', type: 'Timespec', count: 2, marshal: timespecMarshal, unmarshal: timespecUnmarshal },
        { name: 'mtime', type: 'Timespec', count: 2, marshal: timespecMarshal, unmarshal: timespecUnmarshal },
        { name: 'ctime', type: 'Timespec', count: 2, marshal: timespecMarshal, unmarshal: timespecUnmarshal },
        { name: 'X__unused', type: 'int64', count: 3, marshal: nullMarshal, unmarshal: nullUnmarshal, omit: true },
    ],
    alignment: 'natural',
    length: 144,
};
(function (DT) {
    DT[DT["UNKNOWN"] = 0] = "UNKNOWN";
    DT[DT["FIFO"] = 1] = "FIFO";
    DT[DT["CHR"] = 2] = "CHR";
    DT[DT["DIR"] = 4] = "DIR";
    DT[DT["BLK"] = 6] = "BLK";
    DT[DT["REG"] = 8] = "REG";
    DT[DT["LNK"] = 10] = "LNK";
    DT[DT["SOCK"] = 12] = "SOCK";
    DT[DT["WHT"] = 14] = "WHT";
})(exports.DT || (exports.DT = {}));
var DT = exports.DT;
;
var Dirent = (function () {
    function Dirent(ino, type, name) {
        this.ino = ino;
        this.type = type;
        this.name = name;
    }
    Object.defineProperty(Dirent.prototype, "off", {
        get: function () {
            return 0;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(Dirent.prototype, "reclen", {
        get: function () {
            var slen = utf8.utf8ToBytes(this.name).length;
            var nZeros = nzeros(slen);
            return slen + nZeros + 1 + 2 + 8 + 8;
        },
        enumerable: true,
        configurable: true
    });
    return Dirent;
}());
exports.Dirent = Dirent;
function nzeros(nBytes) {
    return (8 - ((nBytes + 3) % 8));
}
function nameMarshal(dst, off, src) {
    if (typeof src !== 'string')
        return [undefined, new Error('src not a string: ' + src)];
    var bytes = utf8.utf8ToBytes(src);
    var nZeros = nzeros(bytes.length);
    if (off + bytes.length + nZeros > dst.byteLength)
        return [undefined, new Error('dst not big enough')];
    for (var i = 0; i < bytes.length; i++)
        dst.setUint8(off + i, bytes[i]);
    for (var i = 0; i < nZeros; i++)
        dst.setUint8(off + bytes.length + i, 0);
    return [bytes.length + nZeros, null];
}
;
function nameUnmarshal(src, off) {
    var len = 0;
    for (var i = off; i < src.byteLength && src.getUint8(i) !== 0; i++)
        len++;
    var str = utf8.utf8Slice(src, off, off + len);
    var nZeros = nzeros(len);
    return [str, len + nZeros, null];
}
;
exports.DirentDef = {
    fields: [
        { name: 'ino', type: 'uint64' },
        { name: 'off', type: 'int64' },
        { name: 'reclen', type: 'uint16' },
        { name: 'type', type: 'uint8' },
        { name: 'name', type: 'string', marshal: nameMarshal, unmarshal: nameUnmarshal },
    ],
    alignment: 'natural',
};

},{"./marshal":7,"./utf8":9}],6:[function(require,module,exports){
'use strict';
var _m = require('./marshal');
var _so = require('./socket');
var _fs = require('./fs');
exports.Marshal = _m.Marshal;
exports.Unmarshal = _m.Unmarshal;
exports.socket = _so;
exports.fs = _fs;

},{"./fs":5,"./marshal":7,"./socket":8}],7:[function(require,module,exports){
'use strict';
function typeLen(type) {
    switch (type) {
        case 'uint8':
        case 'int8':
            return 1;
        case 'uint16':
        case 'int16':
            return 2;
        case 'uint32':
        case 'int32':
            return 4;
        case 'uint64':
        case 'int64':
            return 8;
        case 'float32':
            return 4;
        case 'float64':
            return 8;
        default:
            console.log('unknown type');
            console.log(type);
            debugger;
            return undefined;
    }
}
function fieldLen(field) {
    var len = typeLen(field.type);
    return len * (field.count ? field.count : 1);
}
var WRITE_FNS = {
    uint8: function (buf, off, field) {
        field = field >>> 0;
        buf.setUint8(off, field);
        return [1, null];
    },
    uint16: function (buf, off, field) {
        field = field >>> 0;
        buf.setUint16(off, field, true);
        return [2, null];
    },
    uint32: function (buf, off, field) {
        field = field >>> 0;
        buf.setUint32(off, field, true);
        return [4, null];
    },
    uint64: function (buf, off, field) {
        var lo = field >>> 0;
        var hi = (field / ((-1 >>> 0) + 1)) >>> 0;
        buf.setUint32(off, lo, true);
        buf.setUint32(off + 4, hi, true);
        return [8, null];
    },
    int8: function (buf, off, field) {
        field = field | 0;
        buf.setInt8(off, field);
        return [1, null];
    },
    int16: function (buf, off, field) {
        field = field | 0;
        buf.setInt16(off, field, true);
        return [2, null];
    },
    int32: function (buf, off, field) {
        field = field | 0;
        buf.setInt32(off, field, true);
        return [4, null];
    },
    int64: function (buf, off, field) {
        var lo = field | 0;
        var hi = (field / ((-1 >>> 0) + 1)) | 0;
        buf.setInt32(off, lo, true);
        buf.setInt32(off + 4, hi, true);
        return [8, null];
    },
};
var READ_FNS = {
    uint8: function (buf, off) {
        var field = buf.getUint8(off) >>> 0;
        return [field, 1, null];
    },
    uint16: function (buf, off) {
        var field = buf.getUint16(off, true) >>> 0;
        return [field, 2, null];
    },
    uint32: function (buf, off) {
        var field = buf.getUint32(off, true) >>> 0;
        return [field, 4, null];
    },
    uint64: function (buf, off) {
        var lo = buf.getUint32(off, true);
        var hi = buf.getUint32(off + 4, true);
        if (hi !== 0)
            hi *= ((-1 >>> 0) + 1);
        return [lo + hi, 8, null];
    },
    int8: function (buf, off) {
        var field = buf.getInt8(off) | 0;
        return [field, 1, null];
    },
    int16: function (buf, off) {
        var field = buf.getInt16(off, true) | 0;
        return [field, 2, null];
    },
    int32: function (buf, off) {
        var field = buf.getInt32(off, true) | 0;
        return [field, 4, null];
    },
    int64: function (buf, off) {
        var lo = buf.getInt32(off, true);
        var hi = buf.getInt32(off + 4, true);
        if (hi !== 0)
            hi *= ((-1 >>> 0) + 1);
        return [lo + hi, 8, null];
    },
};
function Marshal(buf, off, obj, def) {
    if (!buf || !obj || !def)
        return [0, new Error('missing required inputs')];
    var start = off;
    var write = WRITE_FNS;
    for (var i = 0; i < def.fields.length; i++) {
        var field = def.fields[i];
        var val = obj[field.name];
        var len = void 0;
        var err = void 0;
        if (field.marshal)
            _a = field.marshal(buf, off, val), len = _a[0], err = _a[1];
        else
            _b = write[field.type](buf, off, val), len = _b[0], err = _b[1];
        if (err)
            return [off - start, err];
        if (len === undefined)
            len = fieldLen(field);
        off += len;
    }
    return [off - start, null];
    var _a, _b;
}
exports.Marshal = Marshal;
function Unmarshal(obj, buf, off, def) {
    if (!buf || !def)
        return [0, new Error('missing required inputs')];
    var start = off;
    var read = READ_FNS;
    for (var i = 0; i < def.fields.length; i++) {
        var field = def.fields[i];
        var val = void 0;
        var len = void 0;
        var err = void 0;
        if (field.unmarshal)
            _a = field.unmarshal(buf, off), val = _a[0], len = _a[1], err = _a[2];
        else
            _b = read[field.type](buf, off), val = _b[0], len = _b[1], err = _b[2];
        if (err)
            return [off - start, err];
        if (!field.omit)
            obj[field.name] = val;
        if (len === undefined)
            len = fieldLen(field);
        off += len;
    }
    return [off - start, null];
    var _a, _b;
}
exports.Unmarshal = Unmarshal;
function isZero(field) {
    for (var i = 0; i < field.byteLength; i++) {
        if (field.getUint8(i) !== 0)
            return false;
    }
    return true;
}
exports.isZero = isZero;

},{}],8:[function(require,module,exports){
'use strict';
var marshal_1 = require('./marshal');
function IPv4BytesToStr(src, off) {
    if (!off)
        off = 0;
    return [
        '' + src.getUint8(off + 0) +
            '.' + src.getUint8(off + 1) +
            '.' + src.getUint8(off + 2) +
            '.' + src.getUint8(off + 3),
        4,
        null
    ];
}
exports.IPv4BytesToStr = IPv4BytesToStr;
function IPv4StrToBytes(dst, off, src) {
    if (!dst || dst.byteLength < 4)
        return [undefined, new Error('invalid dst')];
    dst.setUint8(off + 0, 0);
    dst.setUint8(off + 1, 0);
    dst.setUint8(off + 2, 0);
    dst.setUint8(off + 3, 0);
    var start = 0;
    var n = off;
    for (var i = 0; i < src.length && n < off + 4; i++) {
        if (src[i] === '.') {
            n++;
            continue;
        }
        dst.setUint8(n, dst.getUint8(n) * 10 + parseInt(src[i], 10));
    }
    return [4, null];
}
exports.IPv4StrToBytes = IPv4StrToBytes;
exports.SockAddrInDef = {
    fields: [
        { name: 'family', type: 'uint16' },
        { name: 'port', type: 'uint16' },
        {
            name: 'addr',
            type: 'uint8',
            count: 4,
            JSONType: 'string',
            marshal: IPv4StrToBytes,
            unmarshal: IPv4BytesToStr,
        },
        {
            name: 'zero',
            type: 'uint8',
            count: 8,
            ensure: marshal_1.isZero,
            omit: true,
        },
    ],
    alignment: 'natural',
    length: 16,
};

},{"./marshal":7}],9:[function(require,module,exports){
'use strict';
exports.kMaxLength = 0x3fffffff;
function blitBuffer(src, dst, offset, length) {
    var i;
    for (i = 0; i < length; i++) {
        if ((i + offset >= dst.length) || (i >= src.length))
            break;
        dst[i + offset] = src[i];
    }
    return i;
}
function utf8Slice(buf, start, end) {
    end = Math.min(buf.byteLength, end);
    var res = [];
    var i = start;
    while (i < end) {
        var firstByte = buf.getUint8(i);
        var codePoint = null;
        var bytesPerSequence = (firstByte > 0xEF) ? 4
            : (firstByte > 0xDF) ? 3
                : (firstByte > 0xBF) ? 2
                    : 1;
        if (i + bytesPerSequence <= end) {
            var secondByte = void 0, thirdByte = void 0, fourthByte = void 0, tempCodePoint = void 0;
            switch (bytesPerSequence) {
                case 1:
                    if (firstByte < 0x80) {
                        codePoint = firstByte;
                    }
                    break;
                case 2:
                    secondByte = buf.getUint8(i + 1);
                    if ((secondByte & 0xC0) === 0x80) {
                        tempCodePoint = (firstByte & 0x1F) << 0x6 | (secondByte & 0x3F);
                        if (tempCodePoint > 0x7F) {
                            codePoint = tempCodePoint;
                        }
                    }
                    break;
                case 3:
                    secondByte = buf.getUint8(i + 1);
                    thirdByte = buf.getUint8(i + 2);
                    if ((secondByte & 0xC0) === 0x80 && (thirdByte & 0xC0) === 0x80) {
                        tempCodePoint = (firstByte & 0xF) << 0xC | (secondByte & 0x3F) << 0x6 | (thirdByte & 0x3F);
                        if (tempCodePoint > 0x7FF && (tempCodePoint < 0xD800 || tempCodePoint > 0xDFFF)) {
                            codePoint = tempCodePoint;
                        }
                    }
                    break;
                case 4:
                    secondByte = buf.getUint8(i + 1);
                    thirdByte = buf.getUint8(i + 2);
                    fourthByte = buf.getUint8(i + 3);
                    if ((secondByte & 0xC0) === 0x80 && (thirdByte & 0xC0) === 0x80 && (fourthByte & 0xC0) === 0x80) {
                        tempCodePoint = (firstByte & 0xF) << 0x12 | (secondByte & 0x3F) << 0xC | (thirdByte & 0x3F) << 0x6 | (fourthByte & 0x3F);
                        if (tempCodePoint > 0xFFFF && tempCodePoint < 0x110000) {
                            codePoint = tempCodePoint;
                        }
                    }
            }
        }
        if (codePoint === null) {
            codePoint = 0xFFFD;
            bytesPerSequence = 1;
        }
        else if (codePoint > 0xFFFF) {
            codePoint -= 0x10000;
            res.push(codePoint >>> 10 & 0x3FF | 0xD800);
            codePoint = 0xDC00 | codePoint & 0x3FF;
        }
        res.push(codePoint);
        i += bytesPerSequence;
    }
    return decodeCodePointsArray(res);
}
exports.utf8Slice = utf8Slice;
var MAX_ARGUMENTS_LENGTH = 0x1000;
function decodeCodePointsArray(codePoints) {
    var len = codePoints.length;
    if (len <= MAX_ARGUMENTS_LENGTH) {
        return String.fromCharCode.apply(String, codePoints);
    }
    var res = '';
    var i = 0;
    while (i < len) {
        res += String.fromCharCode.apply(String, codePoints.slice(i, i += MAX_ARGUMENTS_LENGTH));
    }
    return res;
}
function utf8ToBytes(string, units) {
    units = units || Infinity;
    var codePoint;
    var length = string.length;
    var leadSurrogate = null;
    var bytes = [];
    for (var i = 0; i < length; i++) {
        codePoint = string.charCodeAt(i);
        if (codePoint > 0xD7FF && codePoint < 0xE000) {
            if (!leadSurrogate) {
                if (codePoint > 0xDBFF) {
                    if ((units -= 3) > -1)
                        bytes.push(0xEF, 0xBF, 0xBD);
                    continue;
                }
                else if (i + 1 === length) {
                    if ((units -= 3) > -1)
                        bytes.push(0xEF, 0xBF, 0xBD);
                    continue;
                }
                leadSurrogate = codePoint;
                continue;
            }
            if (codePoint < 0xDC00) {
                if ((units -= 3) > -1)
                    bytes.push(0xEF, 0xBF, 0xBD);
                leadSurrogate = codePoint;
                continue;
            }
            codePoint = leadSurrogate - 0xD800 << 10 | codePoint - 0xDC00 | 0x10000;
        }
        else if (leadSurrogate) {
            if ((units -= 3) > -1)
                bytes.push(0xEF, 0xBF, 0xBD);
        }
        leadSurrogate = null;
        if (codePoint < 0x80) {
            if ((units -= 1) < 0)
                break;
            bytes.push(codePoint);
        }
        else if (codePoint < 0x800) {
            if ((units -= 2) < 0)
                break;
            bytes.push(codePoint >> 0x6 | 0xC0, codePoint & 0x3F | 0x80);
        }
        else if (codePoint < 0x10000) {
            if ((units -= 3) < 0)
                break;
            bytes.push(codePoint >> 0xC | 0xE0, codePoint >> 0x6 & 0x3F | 0x80, codePoint & 0x3F | 0x80);
        }
        else if (codePoint < 0x110000) {
            if ((units -= 4) < 0)
                break;
            bytes.push(codePoint >> 0x12 | 0xF0, codePoint >> 0xC & 0x3F | 0x80, codePoint >> 0x6 & 0x3F | 0x80, codePoint & 0x3F | 0x80);
        }
        else {
            throw new Error('Invalid code point');
        }
    }
    return bytes;
}
exports.utf8ToBytes = utf8ToBytes;

},{}]},{},[3]);
`
