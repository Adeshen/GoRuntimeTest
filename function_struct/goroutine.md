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





## 汇编层面调用解读

``` go
package main

import "fmt"

func hello() {

}

func main() {
	go func(i int) {
		fmt.Println(i)
	}(1000)

	go hello()
	stop := make(chan int)
	<-stop

	// runtime.Tcase = impl.H{}
}

```



### main.main

```go
TEXT main.main(SB) D:/SUSYstudy/go_cxl/GoRuntimeTest/function_struct/function.go
  function.go:9		0x48e620		493b6610		CMPQ 0x10(R14), SP			
  function.go:9		0x48e624		7669			JBE 0x48e68f				
  function.go:9		0x48e626		4883ec18		SUBQ $0x18, SP				
  function.go:9		0x48e62a		48896c2410		MOVQ BP, 0x10(SP)			
  function.go:9		0x48e62f		488d6c2410		LEAQ 0x10(SP), BP			
  function.go:10	0x48e634		488d0565da0000		LEAQ type.*+53408(SB), AX		
  function.go:10	0x48e63b		0f1f440000		NOPL 0(AX)(AX*1)			
  function.go:10	0x48e640		e85bdef7ff		CALL runtime.newobject(SB)		
  function.go:10	0x48e645		488d0d54000000		LEAQ main.main.func2(SB), CX		
  function.go:10	0x48e64c		488908			MOVQ CX, 0(AX)				
  function.go:10	0x48e64f		488d0d22390200		LEAQ go.func.*+632(SB), CX		
  function.go:10	0x48e656		48894808		MOVQ CX, 0x8(AX)			
  function.go:10	0x48e65a		e82110fbff		CALL runtime.newproc(SB)		
  function.go:14	0x48e65f		488d050a390200		LEAQ go.func.*+624(SB), AX		
  function.go:14	0x48e666		e81510fbff		CALL runtime.newproc(SB)		
  function.go:15	0x48e66b		488d05ce7b0000		LEAQ type.*+29248(SB), AX		
  function.go:15	0x48e672		31db			XORL BX, BX				
  function.go:15	0x48e674		e86761f7ff		CALL runtime.makechan(SB)		
  function.go:16	0x48e679		31db			XORL BX, BX				
  function.go:16	0x48e67b		0f1f440000		NOPL 0(AX)(AX*1)			
  function.go:16	0x48e680		e8bb70f7ff		CALL runtime.chanrecv1(SB)		
  function.go:19	0x48e685		488b6c2410		MOVQ 0x10(SP), BP			
  function.go:19	0x48e68a		4883c418		ADDQ $0x18, SP				
  function.go:19	0x48e68e		c3			RET					
  function.go:9		0x48e68f		e80cdafcff		CALL runtime.morestack_noctxt.abi0(SB)	
  function.go:9		0x48e694		eb8a			JMP main.main(SB)	
```

0. NOPL 0(AX)(AX*1)

在x86汇编中，`NOPL`是一种无操作指令，它不会对寄存器或内存做出任何改变。该指令通常用于在指令序列中占位或填充空闲的字节，以确保其他指令能够正确对齐。

`NOPL 0(AX)(AX*1)`指令中，`0(AX)`表示在地址`AX`的偏移量为0的位置读取或写入一个字节，`(AX*1)`表示使用寄存器`AX`的值乘以1作为地址的偏移量。因此，该指令是无效的，因为它没有有效的操作数。

0. **x86    CALL runtime.newobject(SB)	  函数返回值会存在哪里？**

在x86汇编中，函数返回值通常存储在寄存器eax中。但是，在这种情况下，函数`runtime.newobject(SB)`并没有返回任何值。相反，它会分配新的内存对象，并返回一个指向该对象的指针，该指针通常存储在寄存器eax中，以便后续代码可以使用该指针来操作该对象。

eax 此时和获取一个栈空间

2. **LEAQ main.main.func2(SB), CX  ; 将函数地址加载到寄存器CX中**

``` asm
LEAQ main.main.func2(SB), CX  ; 将函数地址加载到寄存器CX中
MOVQ CX, 0(AX)                ; 将寄存器CX中的地址推入堆栈
CALL other_function           ; 调用其他函数
```

3. **LEAQ go.func.*+632(SB), CX**

具体来说，`go.func.*+632(SB)`是一个带有通配符的符号引用，它表示某个以`go.func.`为前缀，且距离该符号632字节的函数。该符号通常是由编译器生成的，用于在程序运行时进行函数调用

4. **MOVQ CX, 0x8(AX)**  ; **将go.func.*+632(SB)函数地址加载到寄存器CX中**

5. **CALL runtime.newproc(SB)  调用函数**



### main.main.func2

``` asm
TEXT main.main.func2(SB) D:/SUSYstudy/go_cxl/GoRuntimeTest/function_struct/function.go
  function.go:10	0x48e6a0		493b6610		CMPQ 0x10(R14), SP		
  function.go:10	0x48e6a4		762f			JBE 0x48e6d5			
  function.go:10	0x48e6a6		4883ec10		SUBQ $0x10, SP			
  function.go:10	0x48e6aa		48896c2408		MOVQ BP, 0x8(SP)		
  function.go:10	0x48e6af		488d6c2408		LEAQ 0x8(SP), BP		
  function.go:10	0x48e6b4		4d8b6620		MOVQ 0x20(R14), R12		
  function.go:10	0x48e6b8		4d85e4			TESTQ R12, R12			
  function.go:10	0x48e6bb		751f			JNE 0x48e6dc			
  function.go:10	0x48e6bd		488b5208		MOVQ 0x8(DX), DX		
  function.go:12	0x48e6c1		488b0a			MOVQ 0(DX), CX			
  function.go:12	0x48e6c4		b8e8030000		MOVL $0x3e8, AX			
  function.go:12	0x48e6c9		ffd1			CALL CX				
  function.go:12	0x48e6cb		488b6c2408		MOVQ 0x8(SP), BP		
  function.go:12	0x48e6d0		4883c410		ADDQ $0x10, SP			
  function.go:12	0x48e6d4		c3			RET				
  function.go:10	0x48e6d5		e826d9fcff		CALL runtime.morestack.abi0(SB)	
  function.go:10	0x48e6da		ebc4			JMP main.main.func2(SB)		
  function.go:10	0x48e6dc		4c8d6c2418		LEAQ 0x18(SP), R13		
  function.go:10	0x48e6e1		4d392c24		CMPQ R13, 0(R12)		
  function.go:10	0x48e6e5		75d6			JNE 0x48e6bd			
  function.go:10	0x48e6e7		49892424		MOVQ SP, 0(R12)			
  function.go:10	0x48e6eb		ebd0			JMP 0x48e6bd			

```

1. CMPQ 0x10(R14), SP 

在x86汇编中，`CMPQ`指令用于比较两个操作数的值，并将比较结果存储在标志寄存器中。`CMPQ`指令是带符号比较指令，它将两个操作数视为有符号整数。

`CMPQ 0x10(R14), SP`指令中，`0x10(R14)`表示在寄存器`R14`的值上加上偏移量`0x10`得到一个内存地址，读取该内存地址中的值作为第一个操作数；`SP`表示将栈指针寄存器中的值作为第二个操作数。该指令将这两个操作数进行比较，并将比较结果存储在标志寄存器中。

2. JBE 0x48e6d5

在x86汇编中，`JBE`指令用于执行条件跳转。`JBE`表示"跳转如果以下两种情况之一成立：无符号整数或无符号条件下的比较结果小于或等于零，或有符号条件下的比较结果小于零"。在跳转时，指令使用一个相对于当前指令位置的偏移量作为跳转地址。

`JBE 0x48e6d5`指令中，`0x48e6d5`是跳转的目标地址。具体的偏移量可以在编译时或汇编时计算得到。

在这个例子中，我们使用`CMPQ`指令比较地址中的值和栈指针，然后使用`JBE`指令进行条件跳转。如果栈指针小于等于地址中的值，则跳转到标签`in_bounds`处继续执行程序。否则，程序将继续执行下一条指令。

