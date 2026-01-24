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

	mux.HandleFunc("GET /teachers/{id}/students",handlers.GetStudentByTeacherId)
	mux.HandleFunc("GET /teachers/{id}/studentcount",handlers.GetStudentCountByTeacherId)


	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students/", handlers.AddStudentsHandler)
	mux.HandleFunc("PATCH /students/", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students/", handlers.DeleteStudentsHandler)
	
	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentsHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneStudentsHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteOneStudentHandler)




	
	mux.HandleFunc("GET /execs/", handlers.ExecsHandler)
	return mux
}
