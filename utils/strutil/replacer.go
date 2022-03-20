package strutil

import (
	"reflect"
	"unsafe"
)

//高性能写法，避免内存拷贝
func String2Bytes(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func Byte2String(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

//替换 头部的 old串
func UnsafeReplaceBegin(raw, old, now []byte) []byte {
	var (
		rawlen = len(raw)
		oldlen = len(old)
		nowlen = len(now)
	)
	//log.Printf("%v %v %v", rawlen, oldlen, nowlen)
	var n = rawlen - oldlen + nowlen
	if n <= 0 {
		n = 0
	}
	var newBytes = make([]byte, n)
	for i := 0; i < nowlen; i++ {
		newBytes[i] = now[i]
	}
	////防止数组越界
	//if rawlen > n {
	//	rawlen = n
	//}
	for i := oldlen; nowlen < n; i++ {
		newBytes[nowlen] = raw[i]

		nowlen++
	}
	return newBytes
}
