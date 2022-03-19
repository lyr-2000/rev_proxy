package httpproxy

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_standardUrl(t *testing.T) {
	var p = "http://www.baidu.com/app/url"
	x, _ := url.Parse(p)
	fmt.Println(x.Path)
}
