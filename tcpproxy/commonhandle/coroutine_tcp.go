package commonhandle

import (
	"io"
	"log"
	"net"
	"net/url"
)

//使用 协程处理 tcp 连接

type CommonHandler struct {
}

//开启协程，简历连接
func transferData(con net.Conn, to *url.URL) {
	//defer con.Close()
	target, err := net.Dial("tcp", to.Host)
	if err != nil {
		log.Printf("error occur [%+v]\n", err)
		return
	}
	//defer target.Close()

	go func() {
		io.Copy(con, target)
		con.Close()
		target.Close()
	}()
	go func() {
		io.Copy(target, con)
		target.Close()
		con.Close()
	}()
	//waitExit := make(chan bool, 1)
	//var buf = make([]byte, consts.TcpBufSize)
	//for {
	//	nbuf, err := con.Read(buf)
	//	if err != nil {
	//		log.Printf("unable read  conn data %+v\n", err)
	//		break
	//	}
	//	_, err = target.Write(buf[:nbuf])
	//	if err != nil {
	//		log.Printf("unable write conn data %+v", err)
	//	}
	//}
	//<-waitExit

}
func handle0(from, to *url.URL) error {
	log.Printf("listen host [%v]", from.Host)
	//监听本地端口
	listener, err := net.Listen("tcp", from.Host)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		//处理连接事件
		con, err := listener.Accept()
		if err != nil {
			log.Printf("connct error %+v\n", err)
		}
		go transferData(con, to)

	}

}
func (*CommonHandler) Serve(from, to *url.URL) {
	if from == nil || to == nil {
		panic("empty url config")
	}
	//非阻塞的方式处理多个连接

	go func() {
		err := handle0(from, to)
		if err != nil {
			log.Printf("%+v\n", err)
		}
	}()

}
