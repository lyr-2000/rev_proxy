package errutil

import "runtime"

func ReadStackInfo(buf []byte) []byte {
	if len(buf) == 0 {
		buf = make([]byte, 2048)
	}
	stack := runtime.Stack(buf, false)
	if stack <= 0 {
		return buf[:0]
	}
	return buf[:stack]
}
