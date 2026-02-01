package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func studentRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students/", handlers.AddStudentsHandler)
	mux.HandleFunc("PATCH /students/", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students/", handlers.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentsHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneStudentsHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteOneStudentHandler)
	return mux

}
