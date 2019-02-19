package dtgo

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type MainControler struct {
	// 对应controller和action
	Con, Act string
	// 相应信息
	Rep http.ResponseWriter
	// 参数
	Params map[string]string
	// 路由对应关系
	ConMap map[string]interface{}
	// 公共页面模板
	Layouts string
	// 页面传递参数

}

// 路由
func (mainc *MainControler) Router(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path := url.Path
	if path != "/favicon.ico" {
		// 处理controllre和action
		pathArr := strings.Split(path, "/")
		// 首字母大写
		mainc.Con = strings.Title(pathArr[1])
		mainc.Act = strings.Title(pathArr[2])
		// 处理参数
		for k, v := range url.Query() {
			if len(v) > 0 && v[0] != "" {
				mainc.Params[k] = v[0]
			}
		}
		mainc.Rep = w
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

// 解析页面
func (mainc *MainControler) Display(page string) {
	// tem, _ := template.ParseFiles(page, mainc.Layouts)
	// tem.ExecuteTemplate(mainc.Rep, page)
}
