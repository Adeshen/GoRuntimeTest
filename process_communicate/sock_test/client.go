package main

import (
	"io"
	"log"
	"net"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		return
	}
	println("Client got:", string(buf[0:n]))
}

func main() {
	c, err := net.Dial("unix", "echo.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = c.Write([]byte("hi"))
	if err != nil {
		log.Fatal("write error:", err)
	}
	reader(c)
	time.Sleep(100 * time.Millisecond)
}
