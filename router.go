package dtgo

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/dogntao/dtgo/utils"
)

var (
	//生成的Html保存目录
	htmlOutPath = "./public/static/"
	//静态文件模版目录
	templatePath = "./"
)

// Assign 传递给页面参数
// var Assign = make(map[string]interface{})

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
	// Assign 传递给页面参数
	Assign map[string]interface{}
	// 路由对应关系
	ConMap map[string]interface{}
	// 后台登录
	AdminConMap map[string]interface{}
	// p     sync.Pool
}

// NewRouterStruct 初始化router
func NewRouterStruct() (routerStruct *RouterStruct) {
	routerStruct = &RouterStruct{
		Params:      make(map[string]string),
		Assign:      make(map[string]interface{}),
		ConMap:      make(map[string]interface{}),
		AdminConMap: make(map[string]interface{}),
	}
	// routerStruct.p.New = func() interface{} {
	// 	return &Context{}
	// }
	return
}

// Router 路由方法
func (routerStruct *RouterStruct) Router(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path := url.Path
	fmt.Println("path:", path)
	if strings.Contains(path, "html") {
		t := template.New(path)
		t.Execute(w, nil)
	}
	if path != "/favicon.ico" && !strings.Contains(path, "html") {
		// 处理controllre和action
		pathArr := strings.Split(path, "/")
		// 处理controller和action
		// controller默认为index
		routerStruct.Con = strings.Title("index")
		// method默认为index
		routerStruct.Act = strings.Title("index")
		for k, v := range pathArr {
			// 首字母大写
			v = strings.Title(v)
			if k == 1 && v != "" {
				// controller
				routerStruct.Con = v
			} else if k == 2 && v != "" {
				// action
				routerStruct.Act = v
			}
		}

		// 重置参数
		routerStruct.Params = make(map[string]string)
		// 处理参数
		for k, v := range url.Query() {
			if len(v) > 0 && v[0] != "" {
				routerStruct.Params[k] = v[0]
			}
		}
		routerStruct.Req = r
		routerStruct.Rep = w

		// ctx := routerStruct.p.Get().(*Context)
		// defer routerStruct.p.Put(ctx)
		// ctx.Config(w, r)

		// 优先跳转静态方法
		filePath := routerStruct.getHtmlPath()
		if routerStruct.ExistFile(filePath) {
			http.Redirect(routerStruct.Rep, routerStruct.Req, filePath, http.StatusTemporaryRedirect)
			return
		}
		// 通过反射调用方法
		conv, exist := routerStruct.ConMap[routerStruct.Con]
		if exist == true {
			rv := reflect.ValueOf(conv)
			method := rv.MethodByName(routerStruct.Act)
			if method.IsValid() {
				// 登录验证
				if _, ok := routerStruct.AdminConMap[routerStruct.Con]; ok {
					routerStruct.ALoginCheck()
				}
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

// RegisterAlogin 注册登录controller
func (routerStruct *RouterStruct) RegisterAlogin(cons []string) {
	for _, v := range cons {
		routerStruct.AdminConMap[v] = v
	}
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
	routerStruct.Assign["Con"] = strings.ToLower(routerStruct.Con)
	routerStruct.Assign["Act"] = strings.ToLower(routerStruct.Act)
	err = tem.ExecuteTemplate(routerStruct.Rep, pageName, routerStruct.Assign)
	if err != nil {
		fmt.Println(err)
	}
	// 资讯生成静态文件
	// fmt.Println("con act is", routerStruct.Assign["Con"], routerStruct.Assign["Act"])
	// if routerStruct.Assign["Con"] == "news" && routerStruct.Assign["Act"] == "info" {
	routerStruct.GetGenerateHtml(tem, pageName)
	// }
}

// getHtmlPath 根据url及参数获取静态文件路径
func (routerStruct *RouterStruct) getHtmlPath() (filePath string) {
	con := strings.ToLower(routerStruct.Con)
	act := strings.ToLower(routerStruct.Act)
	fileUrl := templatePath
	// 首页
	if act == "index" {
		fileUrl += con
	}
	if id, ok := routerStruct.Assign["Id"]; ok {
		fileUrl += "/" + id.(string)
	}
	filePath = fileUrl + ".html"
	fmt.Println("con act filePath is", con, act, filePath)
	return
}

// 获取生成静态文件
func (routerStruct *RouterStruct) GetGenerateHtml(tem *template.Template, pageName string) {
	// fileName := pageName
	filePath := routerStruct.getHtmlPath()
	fmt.Println("文件路径", filePath)
	return
	// 1、判断文件是否存在(存在删除)
	if routerStruct.ExistFile(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println(err)
		}
	}
	// 2、打开文件
	fileTmp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("打开文件失败")
	}
	defer fileTmp.Close()

	// 3、生成静态文件
	tem.ExecuteTemplate(fileTmp, pageName, routerStruct.Assign)
	if err != nil {
		fmt.Println(err)
	}
}

//判断文件是否存在
func (routerStruct *RouterStruct) ExistFile(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// 显示后台页面
func (routerStruct *RouterStruct) DisplayAdmin(page string) {
	tem, err := template.ParseFiles(page, "views/layouts/admin/left.html", "views/layouts/admin/header.html", "views/layouts/admin/footer.html")
	if err != nil {
		fmt.Println(err)
	}
	pageName := getFileName(page)
	routerStruct.Assign["Con"] = strings.ToLower(routerStruct.Con)
	routerStruct.Assign["Act"] = strings.ToLower(routerStruct.Act)
	err = tem.ExecuteTemplate(routerStruct.Rep, pageName, routerStruct.Assign)
	if err != nil {
		fmt.Println(err)
	}
}

// 后台登录验证(未登录跳转到登录页)
func (routerStruct *RouterStruct) ALoginCheck() {
	// 获取cookie
	_, err := utils.GetCookie(routerStruct.Req, "user_info")
	if err != nil {
		http.Redirect(routerStruct.Rep, routerStruct.Req, "/admin/index", http.StatusFound)
	}
}
