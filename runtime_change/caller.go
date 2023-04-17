package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lcallee
#include "callee.h"
*/
import "C"

// import (
// 	"fmt"
// )

func main() {
	C.SayHello()
	//fmt.Println("Success!")
}
