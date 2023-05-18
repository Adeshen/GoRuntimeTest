package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

type Args struct {
	A int
	B int
}

type MyService struct {
}

func (*MyService) Multiply(A int, B int) int {
	return A * B
}

func child_node(register_ip string, register_port string, child_port string) {
	myService := new(MyService)
	// 注册RPC服务

	//  hash <uuid, reflect.type>
	rpc.Register(myService)
	// 找到函数

	// 启动RPC服务监听
	listener, err := net.Listen("tcp", ":"+child_port)
	if err != nil {
		fmt.Println("监听端口失败:", err)
		return
	}
	fmt.Println("正在监听端口 " + child_port)

	register_service := func(register_ip string, register_port string) {
		// 待填写

	}
	go register_service(register_ip, register_port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败:", err)
			continue
		}
		go rpc.ServeConn(conn)

		// 根据哈希码找到对应函数并运行
		// 待处理

	}
}

func main_node(register_ip string, register_port string) {

}

func send_request(register_ip string, register_port string) {
	client, err := rpc.Dial("tcp", register_ip+":"+register_port)
	if err != nil {
		fmt.Println("连接RPC服务器失败:", err)
		return
	}

	// 准备RPC请求参数
	args := &Args{A: 5, B: 3}
	var reply int

	// 调用RPC方法
	err = client.Call("MyService.Multiply", args, &reply)
	if err != nil {
		fmt.Println("RPC调用失败:", err)
		return
	}
	fmt.Println("结果:", reply)
}

var node_level bool

func main() {
	// 获取命令行参数
	fmt.Println("命令行参数数量:", len(os.Args))
	for k, v := range os.Args {
		fmt.Printf("args[%v]=[%v]\n", k, v)
	}
	if len(os.Args) >= 5 && os.Args[1] == "child" {
		//子节点
		node_level = false
		child_node(os.Args[2], os.Args[3], os.Args[4])
		return
	} else {
		//注册节点
		node_level = true
		go main_node(os.Args[2], os.Args[3])
	}

	go func() {
		if node_level {
			send_request(os.Args[2], os.Args[3])
			return
		}

	}()

}
