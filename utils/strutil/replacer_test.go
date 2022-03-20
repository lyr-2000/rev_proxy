package strutil

import (
	"fmt"
	"testing"
)

func TestReplaceBegin(t *testing.T) {
	var s = UnsafeReplaceBegin([]byte("/ide/golang/a.png/"), []byte("/ide/golang/"), []byte("/"))
	fmt.Printf("%s\n", string(s))
}
