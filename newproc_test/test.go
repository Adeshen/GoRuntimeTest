package main

import (
	"fmt"
	"runtime"
	"time"
)

func foo(x, y int) (z int) {
	fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
	z = x + y
	for i := 0; i < 10; i++ {
		z++
	}
	// buf := make([]byte, 9024)
	// n := runtime.Stack(buf, false)
	// fmt.Printf("Stack trace:\n%s\n", buf[:n])
	// fmt.Printf("bytes:%d\n", n)
	return
}

type Xcom struct {
	heloo string
	test  int
}

func complex(som Xcom) {
	// fmt.Printf(som.heloo)
	// buf := make([]byte, 9024)
	// n := runtime.Stack(buf, false)
	// fmt.Printf("Stack trace:\n%s\n", buf[:n])
	// fmt.Printf("bytes:%d\n", n)
	printStackAddress()
}

func printStackAddress() {
	// 获取当前协程的栈地址

	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	runtime.CallersFrames(pc[:n])
	frames := runtime.CallersFrames(pc[:n])

	for i := 0; true; i++ {
		frame, more := frames.Next()
		fmt.Printf("Goroutine Stack(%d) file:%s:%d:%s, Address: %v\n",
			i, frame.File, frame.Line, frame.Function, frame.Entry)

		if !more {
			break
		}
	}
}

func main() {
	x := 99
	y := x * x
	z := foo(x, y)

	fmt.Printf("z=%d\n", z)
	go complex(Xcom{heloo: "522223", test: 123})

	// buf := make([]byte, 9024)
	// n := runtime.Stack(buf, false)
	// fmt.Printf("Stack trace:\n%s\n", buf[:n])
	// num := runtime.NumGoroutine()
	// fmt.Printf("Number of goroutines: %d\n", num)

	// p := pprof.Lookup("goroutine")
	// if p != nil {
	// 	p.WriteTo(os.Stdout, 1)
	// 	// profile := p.WriteTo(os.Stdout, 1)
	// 	// fmt.Println(profile)
	// }

	time.Sleep(1 * time.Second)
	// log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
}
