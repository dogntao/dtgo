package dtgo

import (
	"fmt"
	"net/http"
)

type MainControler struct {
	ConMap map[string]interface{}
}

func (mainc *MainControler) Router(w http.ResponseWriter, r *http.Request) {
	fmt.Println(mainc.ConMap)
}
