package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"redis/resp"
	"syscall"
)

var maxBufferSize = 1024

// TIP 通过构建一个简易redis ，来学习redis 底层实现

func main() {
	// 创建一个socket 监听ipv4
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if nil != err {
		panic(err)
	}
	// 关闭文件描述符
	defer syscall.Close(fd)

	// 设置socket 主要是设置 SO_REUSEADDR   让地址可以重复使用,如果服务重启能马上绑定，而不会因为等待旧连接的超时而导致绑定失败。
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); nil != err {
		panic(err)
	}

	//绑定地址和端口
	addr := syscall.SockaddrInet4{Port: 6379, Addr: [4]byte{0, 0, 0, 0}}
	if err = syscall.Bind(fd, &addr); nil != err {
		panic(err)
	}

	//监听端口
	if err = syscall.Listen(fd, syscall.SOMAXCONN); nil != err {
		panic(err)
	}
	fmt.Println("succee list 0.0.0.0 6379 tcp")
	for {
		nfd, _, err := syscall.Accept(fd)
		if nil != err {
			panic(err)
		}
		//  nfd  新的文件描述符
		go handleConnection(nfd)
	}
}
func handleConnection(fd int) {
	defer syscall.Close(fd)
	buf := make([]byte, maxBufferSize)
	file := os.NewFile(uintptr(fd), "socket") // 将文件描述符转换为 *os.File
	if file == nil {
		fmt.Println("Error creating os.File from file descriptor")
		return
	}
	defer file.Close()
	writer := resp.NewWriter(file) // 使用 *os.File 创建一个 Writer
	for {
		n, err := syscall.Read(fd, buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		// 当返回0字节表示已经断开了
		if n == 0 {
			fmt.Println("client disconnected")
			return
		}
		r := resp.NewResp(bytes.NewReader(buf))
		value, err := r.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("unexpected EOF")
				break
			}
			fmt.Println("error reading from client,err = ", err.Error())
			return
		}
		if err := writer.Flush(value); err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}
