package dtgo

import (
	"fmt"
	"net/http"
)

func Router(w http.ResponseWriter, r *http.Request) {
	fmt.Println("123")
	fmt.Println(r)
	fmt.Println(w)
}
