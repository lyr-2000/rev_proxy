package epoll

import (
	"fmt"
	"net"
	"testing"
)

//内置tcpecho 服务器，用于测试连接
func Test_server(t *testing.T) {

	var process = func(conn net.Conn) {
		defer conn.Close()
		for {
			var buf [128]byte
			n, err := conn.Read(buf[:])

			if err != nil {
				fmt.Printf("read from connect failed, err: %v\n", err)
				break
			}
			str := string(buf[:n])
			fmt.Printf("receive from client, data: %v\n", str)
			conn.Write([]byte(fmt.Sprintf("大胆： %s", str)))
		}
	}
	//   simple tcp server
	//1.listen ip+port
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return
	}

	//2.accept client request
	//3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}

		//create goroutine for each connect
		go process(conn)
	}
}
