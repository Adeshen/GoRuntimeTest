package impl

import (
	"fmt"
	e "hook/exec"
	"strconv"
)

type TestImpl struct {
}

func (TestImpl) Test() {
	fmt.Println("Test pass")
}

func (TestImpl) AfterMain() {
	fmt.Println("AfterMain Test pass")
}

func (TestImpl) BeforeMain() {
	print("BeforeMain Test pass\n")
}

func (TestImpl) NewProc(fn uintptr) {
	a, _ := strconv.Atoi(fmt.Sprintf("%d", serveGo))
	if a != int(fn) {
		go serveGo()
	}
}

func serveGo() {
	fmt.Println("serveGo")
}

func init() {
	e.TestVar = TestImpl{}
	fmt.Println("Test init")
}
