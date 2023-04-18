package main

import (
	"fmt"
	"runtime"
)

type H struct {
}

func (H) S() {
	fmt.Println("hello")
}

func main() {
	fmt.Println("server_start")
	runtime.Tcase = H{}

	go func() {

	}()
	// fmt.Print(runtime.runtimeInitTime)
}
