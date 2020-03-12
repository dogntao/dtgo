package dtgo

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"text/template"
)

// Assign 传递给页面参数
var Assign = make(map[string]interface{})

// RouterStruct 路由参数
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

// Router 路由方法
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
		routerStruct.Req = r
		routerStruct.Rep = w

		// method默认为index
		if routerStruct.Act == "" {
			routerStruct.Act = strings.Title("index")
		}

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
	}
}

// RegisterRouter 注册controller文件到路由
func (routerStruct *RouterStruct) RegisterRouter(con string, inter interface{}) {
	routerStruct.ConMap[con] = inter
}

// 根据路径获取文件名(不带后缀)
func getFileName(filePath string) string {
	pageArr := strings.Split(filePath, "/")
	pageName := pageArr[len(pageArr)-1]
	pageNameArr := strings.Split(pageName, ".")
	return pageNameArr[0]
}

// Display 渲染页面方法
func (routerStruct *RouterStruct) Display(page string) {
	tem, err := template.ParseFiles(page, "views/layouts/index/about_left.html", "views/layouts/index/have_left.html", "views/layouts/index/no_left.html", "views/layouts/index/header.html", "views/layouts/index/footer.html")
	if err != nil {
		fmt.Println(err)
	}
	pageName := getFileName(page)
	Assign["Con"] = strings.ToLower(routerStruct.Con)
	Assign["Act"] = strings.ToLower(routerStruct.Act)
	err = tem.ExecuteTemplate(routerStruct.Rep, pageName, Assign)
	if err != nil {
		fmt.Println(err)
	}
}
