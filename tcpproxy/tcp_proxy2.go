package tcpproxy

import (
	"github.com/spf13/viper"
	"log"
)

func _reg(f string, to []interface{}) {
	if len(to) <= 0 {

		return
	}
	//var s = make([]string, len(to))
	//for i, _ := range to {
	//	s[i] = to[i].(string)
	//}
	log.Printf("init tcp proxy from = [%v], to = [%v]\n", f, to)

	err := registerTcpConf(f, to[0].(string))
	if err != nil {
		panic(err)
	}
}

func WithConfig() {
	_array := viper.Get("proxy.tcp")

	arr, isArr := _array.([]interface{})
	if !isArr {
		return
	}

	for _, o := range arr {
		log.Printf("%T\n", o)
		mp, ismp := o.(map[interface{}]interface{})
		if !ismp {

			panic("proxy.http is not array")
		}
		if s, ok := mp["nodes"]; !ok {
			log.Printf("v = %v", s)
		}
		m, ok := mp["nodes"].(map[interface{}]interface{})
		if !ok {
			panic("nodes is not object")
		}
		//is mp
		s := m["from"].(string)
		t := m["to"].([]interface{})
		_reg(s, t)
	}

}
