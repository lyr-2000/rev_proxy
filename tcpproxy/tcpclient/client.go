package tcpclient

import (
	"encoding/json"
	"errors"
	"net"
)

//type TcpClient struct {
//}

func Listen(host string) (net.Listener, error) {
	listener, err := net.Listen("tcp", host)
	//if err != nil {
	//	return nil,err
	//}
	return listener, err
}

func NewClient(host string) (net.Conn, error) {
	c, err := net.Dial("tcp", host)
	return c, err
}

func TcpSend(con net.Conn, msg interface{}) (int, error) {
	if msg == nil {
		return 0, errors.New("empty msg")
	}
	switch msg.(type) {
	case string:
		return con.Write([]byte(msg.(string)))
	case []byte:
		return con.Write(msg.([]byte))
	default:
		bs, err := json.Marshal(msg)
		if err != nil {
			return 0, err
		}
		return con.Write(bs)

	}

}
