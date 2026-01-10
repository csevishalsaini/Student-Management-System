package handlers

import (
	"fmt"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Students route "))
		fmt.Println("Hello GET Method on Students route ")

	case http.MethodPost:
		w.Write([]byte("Hello Post Method on Students route "))
		fmt.Println("Hello Post Method on Students route ")

	case http.MethodPut:
		w.Write([]byte("Hello Put Method on Students route "))
		fmt.Println("Hello Put Method on Students route ")

	case http.MethodPatch:
		w.Write([]byte("Hello Patch Method on Students route "))
		fmt.Println("Hello Patch Method on Students route ")

	case http.MethodDelete:
		w.Write([]byte("Hello Delete Method on Students route "))
		fmt.Println("Hello Delete Method on Students route ")
	}

	w.Write([]byte("Hello students route "))
	fmt.Println("Hello students route ")
}
