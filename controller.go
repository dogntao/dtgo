package dtgo

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type MainControler struct {
	Con, Act string
	Params   map[string]string
	ConMap   map[string]interface{}
}

func (mainc *MainControler) Router(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path := url.Path
	if path != "/favicon.ico" {
		// 处理controllre和action
		pathArr := strings.Split(path, "/")
		mainc.Con = strings.Title(pathArr[1])
		mainc.Act = strings.Title(pathArr[2])
		// 处理参数
		for k, v := range url.Query() {
			if len(v) > 0 && v[0] != "" {
				mainc.Params[k] = v[0]
			}
		}
		fmt.Println(mainc.Params)

		// 通过反射调用方法
		conv, exist := mainc.ConMap[mainc.Con]
		if exist == true {
			rv := reflect.ValueOf(conv)
			method := rv.MethodByName(mainc.Act)
			if method.IsValid() {
				method.Call([]reflect.Value{})
			} else {
				fmt.Println("method is not exist")
			}
		} else {
			fmt.Println("controller is not exist")
		}
	}
	// fmt.Println(url.RawQuery)
}
