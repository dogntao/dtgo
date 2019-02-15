package dtgo

import "net/http"

func Run() {
	http.HandleFunc("/", Router)
	http.ListenAndServe("127.0.0.1", nil)
}
