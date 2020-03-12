package dtgo

import "net/http"

// Host 初始化web服务参数
type Host struct {
	Public string // 静态路径名称(/public)
	Js     string // 静态js路径(/js/)
	CSS    string // 静态css路径(/css/)
	Image  string // 静态图片路径(/image/)
	Port   string // 端口
	Router *RouterStruct
}

// NewHost 初始化Host
func NewHost(public, js, css, image, port string) *Host {
	return &Host{
		Public: public,
		Js:     js,
		CSS:    css,
		Image:  image,
		Port:   port,
		Router: &RouterStruct{
			Params: make(map[string]string),
			ConMap: make(map[string]interface{}),
		},
	}
}

// Run 启动服务
func (host *Host) Run() {
	// js路径
	http.Handle(host.Js, http.FileServer(http.Dir(host.Public)))
	// css路径
	http.Handle(host.CSS, http.FileServer(http.Dir(host.Public)))
	// 静态图片路径
	http.Handle(host.Image, http.FileServer(http.Dir(host.Public)))
	http.HandleFunc("/", host.Router.Router)
	http.ListenAndServe(host.Port, nil)
}
