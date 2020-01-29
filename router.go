package dtgo

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type RouterStruct struct {
	// 对应controller和action
	Con, Act string
	// 请求信息
	Req *http.Request
	// 返回信息
	Rep http.ResponseWriter
	// 参数
	Params map[string]string
	// 路由对应关系
	ConMap map[string]interface{}
}

// 路由
func (routerStruct *RouterStruct) Router(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path := url.Path
	if path != "/favicon.ico" {
		// 处理controllre和action
		pathArr := strings.Split(path, "/")
		// 处理controller和action
		for k, v := range pathArr {
			// 首字母大写
			v = strings.Title(v)
			if k == 1 {
				// controller
				routerStruct.Con = v
			} else if k == 2 {
				// action
				routerStruct.Act = v
			}
		}
		// 处理参数
		routerStruct.Params = make(map[string]string)
		for k, v := range url.Query() {
			if len(v) > 0 && v[0] != "" {
				routerStruct.Params[k] = v[0]
			}
		}
		routerStruct.Rep = w
		// 通过反射调用方法
		conv, exist := routerStruct.ConMap[routerStruct.Con]
		if exist == true {
			rv := reflect.ValueOf(conv)
			method := rv.MethodByName(routerStruct.Act)
			if method.IsValid() {
				method.Call([]reflect.Value{})
			} else {
				fmt.Println(routerStruct.Act + " method is not exist")
			}
		} else {
			fmt.Println(routerStruct.Con + " controller is not exist")
		}
		fmt.Println(routerStruct)
	}
}

// 注册路由
func (routerStruct *RouterStruct) RegisterRouter(con string, inter interface{}) {
	routerStruct.ConMap[con] = inter
}
