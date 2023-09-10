 package main

import (
	"fmt"
	"time"
)

func foo(x, y int) (z int) {
        fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
        z = x + y

        return
}

type Xcom struct{
	heloo string
	test int
}

func complex(som Xcom){
	fmt.Printf(som.heloo)
}

func main() {
        x := 99
        y := x * x
        z := foo(x, y)

        fmt.Printf("z=%d\n", z)
		go complex(Xcom{heloo:"123",test:123})

		time.Sleep(1*time.Second)
}