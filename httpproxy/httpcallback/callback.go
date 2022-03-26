package httpcallback

import (
	"net/http"
)

//http异常回调
func OnError(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("%v, %v, %v\n", r.Host, r.Header, err)
	//log.Printf("status code %v", r.Response)

}

//修改回调

func ModifyResponse(r *http.Response) error {
	if r.StatusCode == 200 {
		//log.Printf("请求成功")

	}
	
	//log.Printf("conn = %s\n", r.Header.Get("Connection"))
	//判断状态码，用来记录统计
	return nil
}
