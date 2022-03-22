package loadbalance

import (
	"hash/crc32"
	"myproxyHttp/utils/serverutil"
	"myproxyHttp/utils/strutil"
	"net/http"
)

type IphashHandler struct {
}

func StringHash(key string) uint32 {
	ieee := crc32.ChecksumIEEE(strutil.String2Bytes(&key))
	return ieee
}

func (*IphashHandler) NextIndex(w http.ResponseWriter, r *http.Request, hostCnt int) int {
	if hostCnt <= 0 {
		return 0
	}
	ip := serverutil.ClientIP(r)
	return int(StringHash(ip)) % hostCnt
}

func (*IphashHandler) RemoveHost(key string) {

}
