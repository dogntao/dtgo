package dtgo

import (
	"fmt"
	"net/http"
)

type MainControler struct {
	conMap map[string]interface{}
}

func (main *MainControler) Router(w http.ResponseWriter, r *http.Request) {
	fmt.Println("123")
	fmt.Println(r)
	fmt.Println(w)
}
