package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"strconv"
	"time"
)

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	// studentList := make([]models.Student, 0)
	// firstName := r.URL.Query().Get("first_name")
	// lastName := r.URL.Query().Get("last_name")

	var execs []models.Exec
	execs, err := sqlconnect.GetExecsDbOperation(execs, r)
	// if shouldReturn {
	// 	return
	// }

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execs),
		Data:   execs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	exec, err := sqlconnect.GetExecById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(exec)

}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {

	var newExecs []models.Exec
	var rawExec []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request ", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rawExec)
	if err != nil {
		http.Error(w, "Invalid request body ", http.StatusBadRequest)
		return
	}

	fields := GetFieldsName(models.Exec{})

	allowedField := make(map[string]struct{})
	for _, field := range fields {
		allowedField[field] = struct{}{}
	}

	for _, exec := range rawExec {
		for key := range exec {
			_, ok := allowedField[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields.", http.StatusBadRequest)
				return
			}

		}
	}

	err = json.Unmarshal(body, &newExecs)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)

		return
	}

	for _, exec := range newExecs {
		err := CheckBlankField(exec)
		if err != nil {
			http.Error(w, "Incorrect field error, Check your field", http.StatusBadRequest)
			return
		}
	}

	addedExecs, err := sqlconnect.AddExecsDBHandler(newExecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	json.NewEncoder(w).Encode(&response)

}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := sqlconnect.PatchExecs(updates)
	if err != nil {
		http.Error(w, "Error to connect with database ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func PatchOneExecsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid exec id", http.StatusBadRequest)
		return
	}
	updated := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Invalid payload request", http.StatusInternalServerError)
		return
	}

	existingExec, err := sqlconnect.PatchOneExec(id, updated)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingExec)
}

func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid exec id", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneExec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "Exec successfully deleted ",
		Id:     id,
	}
	json.NewEncoder(w).Encode(response)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required ", http.StatusBadRequest)
	}

	user, err := sqlconnect.GetUserByUserName(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password ", http.StatusForbidden)
		return
	}

	if user.InactiveStatus {
		http.Error(w, "Account is incative", http.StatusForbidden)
		return
	}

	err = utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	tokenString, err := utils.SignToken(user.ID, req.Username, user.Role)

	if err != nil {
		http.Error(w, "Could not create login token ", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "test",
		Value:    "testing",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(response)
}


func LogOutHandler(w http.ResponseWriter, r *http.Request){
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0,0),
		SameSite: http.SameSiteStrictMode,

	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{message: "Logged out Succesfully"}`))
}