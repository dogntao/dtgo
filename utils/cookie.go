package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

// 设置cookie
func SetCooke(rsp http.ResponseWriter, key string, value map[string]string) {
	// map转换成json
	valueByte, _ := json.Marshal(value)
	// josn base64 encode
	valueStr := base64.StdEncoding.EncodeToString(valueByte)
	// set cookie
	cookie := &http.Cookie{Name: key, Value: valueStr, Path: "/", MaxAge: 86400}
	http.SetCookie(rsp, cookie)
}

// 获取cookie
func GetCookie(req *http.Request, key string) (map[string]interface{}, error) {
	cookie, _ := req.Cookie(key)
	if cookie != nil {
		valStr := cookie.Value
		// base64 decode
		valByte, _ := base64.StdEncoding.DecodeString(valStr)
		// json decode
		cookieMap := make(map[string]interface{})
		json.Unmarshal(valByte, &cookieMap)
		return cookieMap, nil
	} else {
		return nil, errors.New("no cookie")
	}
}

// 删除cookie
func DelCookie(rsp http.ResponseWriter, key string) {
	// set cookie
	cookie := &http.Cookie{Name: key, Path: "/", MaxAge: -1}
	http.SetCookie(rsp, cookie)
}
