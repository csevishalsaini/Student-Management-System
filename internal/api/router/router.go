package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.RootHandler)

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeachersHandler)
	
	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeachersHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteOneTeacherHandler)
	



	mux.HandleFunc("/students/", handlers.StudentsHandler)
	mux.HandleFunc("GET /students/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /students/", handlers.AddTeachersHandler)
	mux.HandleFunc("PATCH /students/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /students/", handlers.DeleteTeachersHandler)
	
	mux.HandleFunc("GET /students/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneTeachersHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteOneTeacherHandler)


	
	mux.HandleFunc("/execs/", handlers.ExecsHandler)
	return mux
}
