package strutil

import "encoding/json"

func ToJSON(s interface{}) string {
	marshal, _ := json.Marshal(s)
	return string(marshal)
}
