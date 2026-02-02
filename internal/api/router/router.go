package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/students/", studentRouter())
	mux.Handle("/teachers/", teacherRouter())
	mux.Handle("/execs/", execsRouter())
	return mux

}
