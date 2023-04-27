package rpcimpl

import (
	"fmt"
	"runtime"
)

type H struct {
}

func (H) S() {
	fmt.Println("hello")
}

func init() {
	runtime.Tcase = H{}
}
