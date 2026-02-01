package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/students/", studentRouter())
	mux.Handle("/teachers/", teacherRouter())

	return mux

	// mux := http.NewServeMux()
	// mux.HandleFunc("GET /", handlers.RootHandler)
	// mux.HandleFunc("GET /execs/", handlers.ExecsHandler)

}
