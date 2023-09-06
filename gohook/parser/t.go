package main

func before_main() bool {
	println("in before_main")
	return true
}

func useless() {
	if before_main() {
		println("if before_main")
		return
	}
}

func main() {
	println("hello")
	f := func(a int) {
		print(a)
	}
	go f(4)
}
