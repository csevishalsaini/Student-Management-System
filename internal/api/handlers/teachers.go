package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// teacherList := make([]models.Teacher, 0)
	// firstName := r.URL.Query().Get("first_name")
	// lastName := r.URL.Query().Get("last_name")

	var teachers [] models.Teacher
	teachers, err := sqlconnect.GetTeachersDbOperation(teachers,r)
	// if shouldReturn {
	// 	return
	// }

	if err != nil{
		return
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		fmt.Println(err)
	}

	teacher, err := sqlconnect.GetTeacherById(id)
	if err!=nil{
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(teacher)

}

func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	addedTeachers, err := sqlconnect.AddTeachersDbHandler(newTeachers)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(&response)

}


func UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	fmt.Println(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid payload request ", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err!= nil{
		log.Println(err)
		http.Error(w,"unable to connect with database ",http.StatusInternalServerError)
	}
	fmt.Println(updatedTeacher)

	updatedTeacherFromDb,err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDb)

}



func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	

	var updates []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := sqlconnect.PatchTeachers(updates)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func PatchOneTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}
	updated := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Invalid payload request", http.StatusInternalServerError)
		return
	}

	existingTeacher, err := sqlconnect.PatchOneTeacher(id, updated)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}



func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil{
		http.Error(w, "Invalid payload request", http.StatusInternalServerError)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachers(w, ids)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:status`
		DeletedIds []int
	}{
		Status:     "Teachers successfully",
		DeletedIds: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}

func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneTeacher(id)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "Teacher successfully deleted ",
		Id:     id,
	}
	json.NewEncoder(w).Encode(response)

}

