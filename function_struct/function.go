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
