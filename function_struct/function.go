package main

import "fmt"

func main() {
	go func(i int) {
		fmt.Println(i)
	}(1000)

	stop := make(chan int)
	<-stop

	// runtime.Tcase = impl.H{}
}
