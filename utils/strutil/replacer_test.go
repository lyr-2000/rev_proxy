package strutil

import (
	"fmt"
	"testing"
)

func TestReplaceBegin(t *testing.T) {
	var s = UnsafeReplaceBegin([]byte("/ide/golang/apng/"), []byte("/ide/golang/"), []byte("/axx/"))
	fmt.Printf("%s\n", string(s))
}
