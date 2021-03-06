package tcpproxy

import (
	"net/url"
)

type TcpReverseProxy interface {
	//非阻塞
	Serve(from, to *url.URL)
}

func DoneFinal() {

}

func registerTcpConf(f, t string) error {
	from, err := url.Parse(f)
	if err != nil {
		return err
	}
	to, err := url.Parse(t)

	if err != nil {
		return err
	}
	proxy := Default()
	//log.Printf("start tcp proxy [%v]\n", from)
	//开启代理
	proxy.Serve(from, to)

	return nil
}

//
//
//func ParseConfigTcp(mp map[string]interface{}) error {
//	if mp == nil {
//		log.Println("无http配置")
//		return nil
//	}
//	log.Printf("config map %#v", strutil.ToJSON(mp))
//	var proxy = mp["proxy"].(map[string]interface{})
//	if httpConf, ok := proxy["tcp"]; ok {
//		if !reflectutil.IsArrayOrSlice(httpConf) {
//			// 如果不是数组
//			return errors.New("配置不正确，proxy.http不是数组")
//		}
//		arr := httpConf.([]interface{})
//		for i, _ := range arr {
//			h := arr[i].(map[string]interface{})
//			err := registerTcpConf(h["from"].(string), h["to"].(string))
//			if err != nil {
//				log.Printf("%+v\n", err)
//			}
//		}
//		// proxy parse
//	}
//
//	return nil
//}
