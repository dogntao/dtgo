package dtgo

import "net/http"

var MainC *MainControler

func Run() {
	http.HandleFunc("/", MainC.Router)
	http.ListenAndServe("127.0.0.1", nil)
}
