# Target

to find out something interesting in Go runtime

* [X] insert code into runtime
* [X] add a pipe in runtime
* [ ] 

# 1.insert code into runtime

you can use the file in `runtime_change` to test after inserting code

# 2.process_communicate



## Add a pipe  runtime

/usr/lib/go-1.18/src/runtime/sys_netbsd_386.s:// func pipe2(flags int32) (r, w int32, errno int32)

/usr/lib/go-1.18/src/runtime/syscall_solaris.go

```go
func syscall_write(fd, buf, nbyte uintptr) (n, err uintptr) 


```



/usr/lib/go-1.18/src/os

```go
func Pipe() (r *File, w *File, err error) {
	var p [2]syscall.Handle
	e := syscall.Pipe(p[:])
	if e != nil {
		return nil, nil, NewSyscallError("pipe", e)
	}
	return newFile(p[0], "|0", "pipe"), newFile(p[1], "|1", "pipe"), nil
}
```



/usr/lib/go-1.18/src/syscall/syscall_windows.go

```
func Pipe(p []Handle) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var r, w Handle
	e := CreatePipe(&r, &w, makeInheritSa(), 0)
	if e != nil {
		return e
	}
	p[0] = r
	p[1] = w
	return nil
}
```



/





