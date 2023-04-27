# go 协程创建函数





## 获取G



```go
func newproc1(fn *funcval, callergp *g, callerpc uintptr) *g {
	***
	newg := gfget(_p_)
	***
}

func gfget(_p_ *p) *g {
retry:
    
    //如果当前协程栈为空，同时全局队列中不空，那么在全局队列中获取
	if _p_.gFree.empty() && (!sched.gFree.stack.empty() || !sched.gFree.noStack.empty()) {
		lock(&sched.gFree.lock)
		// Move a batch of free Gs to the P.
		for _p_.gFree.n < 32 {
			// Prefer Gs with stacks.
			gp := sched.gFree.stack.pop()  //在调度器的全局队列中获取
			if gp == nil {
				gp = sched.gFree.noStack.pop()
				if gp == nil {
					break
				}
			}
			sched.gFree.n--
			_p_.gFree.push(gp)
			_p_.gFree.n++
		}
		unlock(&sched.gFree.lock)
		goto retry
	}
	gp := _p_.gFree.pop()  //在当前p中获取
	if gp == nil {
		return nil
	}
	_p_.gFree.n--   
    //清除 原来 老堆栈的空间
	if gp.stack.lo != 0 && gp.stack.hi-gp.stack.lo != uintptr(startingStackSize) {
		// Deallocate old stack. We kept it in gfput because it was the
		// right size when the goroutine was put on the free list, but
		// the right size has changed since then.
		systemstack(func() {
			stackfree(gp.stack)
			gp.stack.lo = 0
			gp.stack.hi = 0
			gp.stackguard0 = 0
		})
	}
	if gp.stack.lo == 0 {
		// Stack was deallocated in gfput or just above. Allocate a new one.
		systemstack(func() {
			gp.stack = stackalloc(startingStackSize)
		})
		gp.stackguard0 = gp.stack.lo + _StackGuard
	} else {
		if raceenabled {
			racemalloc(unsafe.Pointer(gp.stack.lo), gp.stack.hi-gp.stack.lo)
		}
		if msanenabled {
			msanmalloc(unsafe.Pointer(gp.stack.lo), gp.stack.hi-gp.stack.lo)
		}
		if asanenabled {
			asanunpoison(unsafe.Pointer(gp.stack.lo), gp.stack.hi-gp.stack.lo)
		}
	}
	return gp
}

//如果都找不到，那么就分配一个全新的g
// Allocate a new g, with a stack big enough for stacksize bytes.
func malg(stacksize int32) *g {
	newg := new(g)
	if stacksize >= 0 {
		stacksize = round2(_StackSystem + stacksize)
		systemstack(func() {
			newg.stack = stackalloc(uint32(stacksize))
		})
		newg.stackguard0 = newg.stack.lo + _StackGuard
		newg.stackguard1 = ^uintptr(0)
		// Clear the bottom word of the stack. We record g
		// there on gsignal stack during VDSO on ARM and ARM64.
		*(*uintptr)(unsafe.Pointer(newg.stack.lo)) = 0
	}
	return newg
}
```



## 初始化g

```go

func newproc1(){
	//***
    //给一个初始帧大小
	totalSize := uintptr(4*goarch.PtrSize + sys.MinFrameSize) // extra space in case of reads slightly beyond frame
	totalSize = alignUp(totalSize, sys.StackAlign)
	//从高位向低位增长
    sp := newg.stack.hi - totalSize
	spArg := sp
	if usesLR {
		// caller's LR
		*(*uintptr)(unsafe.Pointer(sp)) = 0
		prepGoExitFrame(sp)
		spArg += sys.MinFrameSize
	}
	//
	memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
	newg.sched.sp = sp
	newg.stktopsp = sp
	newg.sched.pc = abi.FuncPCABI0(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
	newg.sched.g = guintptr(unsafe.Pointer(newg))
	gostartcallfn(&newg.sched, fn)
	newg.gopc = callerpc
	newg.ancestors = saveAncestors(callergp)
	newg.startpc = fn.fn
	if isSystemGoroutine(newg, false) {
    //***
}
    

//stack.go
func gostartcallfn(gobuf *gobuf, fv *funcval) {
	var fn unsafe.Pointer
	if fv != nil {
		fn = unsafe.Pointer(fv.fn)
	} else {
		fn = unsafe.Pointer(abi.FuncPCABIInternal(nilfunc))
	}
	gostartcall(gobuf, fn, unsafe.Pointer(fv))
}
    
    // Stack frame layout
//
// (x86)
// +------------------+
// | args from caller |
// +------------------+ <- frame->argp
// |  return address  |
// +------------------+
// |  caller's BP (*) | (*) if framepointer_enabled && varp < sp
// +------------------+ <- frame->varp
// |     locals       |
// +------------------+
// |  args to callee  |
// +------------------+ <- frame->sp

    
//sys_86.go    
func gostartcall(buf *gobuf, fn, ctxt unsafe.Pointer) {
	sp := buf.sp
	sp -= goarch.PtrSize
	*(*uintptr)(unsafe.Pointer(sp)) = buf.pc
	buf.sp = sp
	buf.pc = uintptr(fn)
	buf.ctxt = ctxt
}

```

在这个函数中初始化了

sched：暂时调度状态

​	sp: 栈顶指针    .   第一次设置的是栈顶的最小值，

​	pc：当前执行命令

​	ctxt： 上下文

gopc: 调用者 go关键字的位置

ancestor：新协程的祖先

startpc: 新协程的起点





具体来说，`newg.sched.pc` 是新 Goroutine 的调度器上下文中的一个字段，它表示 Goroutine 要执行的下一条指令的地址。在这里，它被设置为调用 `goexit` 函数的地址。

`abi.FuncPCABI0` 是一个函数，它返回一个函数的入口地址，这里的参数 `goexit` 是 Go 运行时库中的一个函数，用于协程退出时的清理工作。`sys.PCQuantum` 是一个常量，表示机器指令的大小。

在这里加上 `sys.PCQuantum` 是为了保证调用 `goexit` 的指令与新 Goroutine 调用该指令之前的指令在同一个函数中。这样做是为了避免出现栈溢出等问题，因为在不同的函数之间跳转时，调用者的栈帧需要被弹出，新函数的栈帧需要被压入，这样会增加栈空间的使用量，可能会导致栈溢出。

总之，这段代码的作用是为新的 Goroutine 设置调度器上下文中的 pc 字段，使其指向 `goexit` 函数的入口地址，并通过加上 `sys.PCQuantum` 保证调用该函数的指令与新 Goroutine 调用该指令之前的指令在同一个函数中。这样做可以保证 Goroutine 的调度和执行过程正确无误，并避免出现栈溢出等问题。







## 测试获取参数

``` go
	gostartcallfn(&newg.sched, fn)
	newg.gopc = callerpc
	newg.ancestors = saveAncestors(callergp)
	newg.startpc = fn.fn

	// var i unsafe.Pointer = 0
	// for ; i < 10; i++ {
	// 	print(*(*uintptr)(unsafe.Pointer(newg.sched.sp - unsafe.Pointer(i)*goarch.PtrSize)))
	// }
	print(*(*uintptr)(unsafe.Pointer(newg.sched.sp + 2*goarch.PtrSize)))
	if isSystemGoroutine(newg, false) {
		atomic.Xadd(&sched.ngsys, +1)
	}

```

但是发现其中

newg.sched.sp + 0*goarch.PtrSize   有数值，其余都是0













