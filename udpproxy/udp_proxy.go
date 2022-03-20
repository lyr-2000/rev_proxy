package udpproxy

import (
	"errors"
	"fmt"
	"log"
	"myproxyHttp/consts"
	"myproxyHttp/utils/reflectutil"
	"myproxyHttp/utils/strutil"
	"net"
	"net/url"
)

type UdpHandler struct {
}

func (*UdpHandler) Serve(f, t *url.URL) {

	go Serve0(f, t)
}

func Serve0(f, t *url.URL) error {
	udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", f.Hostname(), f.Port()))
	servConn, _ := net.ListenUDP("udp", udpAddr)
	defer servConn.Close()

	//from => to

	target, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", t.Hostname(), t.Port()))
	if err != nil {
		log.Printf("resolve error %+v\n", err)
		return err
	}
	buf := make([]byte, consts.UdpBufSize)
	for {
		//n,remote,err
		n, _, e := servConn.ReadFromUDP(buf)
		if n <= 0 || e != nil {
			continue
		}
		//fmt.Printf("receive pack %s\n", buf[:n])
		n, e = servConn.WriteToUDP(buf[:n], target)
		if e != nil {
			log.Printf("eror udp write %+v\n", e)
		}
		//log.Printf("write n %v\n", n)
	}

}

//------------

func registerUdpConf(f, t string) error {
	from, err := url.Parse(f)
	if err != nil {
		return err
	}
	to, err := url.Parse(t)

	if err != nil {
		return err
	}
	proxy := &UdpHandler{}
	log.Printf("start udp proxy [%v]\n", from)
	//开启代理
	proxy.Serve(from, to)

	return nil
}
func ParseConfigUdp(mp map[string]interface{}) error {
	if mp == nil {
		log.Println("无http配置")
		return nil
	}
	log.Printf("config map %#v", strutil.ToJSON(mp))
	var proxy = mp["proxy"].(map[string]interface{})
	if httpConf, ok := proxy["udp"]; ok {
		if !reflectutil.IsArrayOrSlice(httpConf) {
			// 如果不是数组
			return errors.New("配置不正确，proxy.http不是数组")
		}
		arr := httpConf.([]interface{})
		for i, _ := range arr {
			h := arr[i].(map[string]interface{})
			err := registerUdpConf(h["from"].(string), h["to"].(string))
			if err != nil {
				log.Printf("%+v\n", err)
			}
		}
		// proxy parse
	}

	return nil
}
