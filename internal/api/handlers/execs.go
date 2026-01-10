package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Query ", r.URL.Query())
	fmt.Println("name ", r.URL.Query().Get("name"))
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Execs route "))

	case http.MethodPost:
		w.Write([]byte("Hello Post Method on Execs route "))

	case http.MethodPut:
		w.Write([]byte("Hello Put Method on Execs route "))

	case http.MethodPatch:
		w.Write([]byte("Hello Patch Method on Execs route "))

	case http.MethodDelete:
		w.Write([]byte("Hello Delete Method on Execs route "))
	}

	w.Write([]byte("Hello execs route "))
	fmt.Println("Hello execs route ")
}
