package udpproxy

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func TestParseConfigUdp(t *testing.T) {
	udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:8078")
	servConn, _ := net.ListenUDP("udp", udpAddr)
	defer servConn.Close()

	//from => to

	//target, err := net.ResolveUDPAddr("udp", t.Host)

	buf := make([]byte, 1024)
	for {
		//n,remote,err
		n, _, e := servConn.ReadFromUDP(buf)
		if n <= 0 || e != nil {
			log.Printf("%+v\n", e)
			continue
		}
		fmt.Println(string(buf[:n]))
	}
}
