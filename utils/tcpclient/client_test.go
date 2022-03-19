package tcpclient

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    net.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_receive(t *testing.T) {
	var (
		host = "localhost:8078"
	)
	//go func() {
	listen, _ := Listen(host)
	defer listen.Close()
	for {
		accept, _ := listen.Accept()
		var buf = make([]byte, 10240)
		n, _ := accept.Read(buf)
		fmt.Println(string(buf[:n]))
	}
	//}()
}
func TestTcpSend(t *testing.T) {
	var (
		host = "localhost:8078"
	)
	go func() {
		listen, _ := Listen(host)
		defer listen.Close()
		for {
			accept, _ := listen.Accept()
			var buf = make([]byte, 10240)
			n, _ := accept.Read(buf)
			fmt.Println(string(buf[:n]))
		}
	}()
	client, err := NewClient(host)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	TcpSend(client, "helll world")
}

func Test_tcp_send1111(t *testing.T) {
	var host = ":8080"
	client, err := NewClient(host)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	TcpSend(client, "helll world")
}
