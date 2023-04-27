package main

import (
	"fmt"
	"runtime"
)

type H struct {
}

func (H) S() {
	fmt.Println("hello  ")
}

func main() {
	fmt.Println("server_start")
	runtime.Tcase = H{}

	he := make(chan int)

	go func(i int) {
		fmt.Println(i)
	}(10)
	// fmt.Print(runtime.runtimeInitTime)
	<-he
}
