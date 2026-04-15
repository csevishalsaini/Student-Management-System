package sqlconnect

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/pkg/utils"
	"strconv"
	"strings"
)

func DeleteOneExec(id int) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Unable to connect Databases")
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Error deleting exec ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error deleting exec")
	}
	fmt.Println(result.RowsAffected())
	rowsEffected, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrieving delete result ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error retrieving delete result")
	}

	if rowsEffected == 0 {
		// http.Error(w, "Exec not found ", http.StatusNotFound)
		return utils.ErrorHandler(err, "Exec not found")
	}
	return nil
}

func PatchOneExec(id int, updated map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return models.Exec{}, utils.ErrorHandler(err, "Unable to connect Databases")
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Exec Not found ", http.StatusNotFound)
			return models.Exec{}, utils.ErrorHandler(err, "Exec Not found")
		}
		// http.Error(w, "Unable to retrieve data ", http.StatusInternalServerError)
		return models.Exec{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	for k, v := range updated {
		switch k {
		case "first_name":
			existingExec.FirstName = v.(string)
		case "last_name":
			existingExec.LastName = v.(string)
		case "email":
			existingExec.Email = v.(string)

		}

		execVal := reflect.ValueOf(&existingExec).Elem()
		fmt.Println(execVal, "  ,, ")
		execType := execVal.Type()

		for k, v := range updated {
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)
					if fieldVal.CanSet() {
						fieldVal.Set(
							reflect.ValueOf(v).Convert(fieldVal.Type()),
						)
					}
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ? ", existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.ID)
	if err != nil {
		// http.Error(w, "Error updating exec ", http.StatusInternalServerError)
		return models.Exec{}, utils.ErrorHandler(err, "Error updating exec")
	}
	return existingExec, nil
}

func PatchExecs(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "Unable to connect database")
	}
	defer db.Close()

	trx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Error starting transaction")
	}

	// Safety rollback (only runs if commit not reached)
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	for _, update := range updates {

		// ---- Handle ID safely ----
		var id int
		switch v := update["id"].(type) {
		case string:
			id, err = strconv.Atoi(v)
			if err != nil {
				return utils.ErrorHandler(err, "Invalid exec ID format")
			}
		case float64: // JSON numbers come as float64
			id = int(v)
		default:
			return utils.ErrorHandler(nil, "Invalid exec ID type")
		}

		// ---- Fetch existing record ----
		var exec models.Exec
		err = trx.QueryRow(
			`SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?`,
			id,
		).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username)

		if err == sql.ErrNoRows {
			return utils.ErrorHandler(err, "Exec not found")
		}
		if err != nil {
			return utils.ErrorHandler(err, "Error retrieving exec")
		}

		// ---- Apply PATCH dynamically ----
		execVal := reflect.ValueOf(&exec).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]

				if jsonTag == k {
					fieldVal := execVal.Field(i)

					if !fieldVal.CanSet() {
						continue
					}

					val := reflect.ValueOf(v)

					if val.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(val.Convert(fieldVal.Type()))
					} else {
						return utils.ErrorHandler(nil, "Type mismatch in PATCH")
					}
					break
				}
			}
		}

		// ---- Update DB ----
		_, err = trx.Exec(
			`UPDATE execs 
			 SET first_name=?, last_name=?, email=?, username=?
			 WHERE id=?`,
			exec.FirstName,
			exec.LastName,
			exec.Email,
			exec.Username,
			exec.ID,
		)

		if err != nil {
			return utils.ErrorHandler(err, "Error updating exec")
		}
	}

	// ---- Commit ----
	if err = trx.Commit(); err != nil {
		return utils.ErrorHandler(err, "Error committing transaction")
	}

	return nil
}

func AddExecsDBHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connect error")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))
	for i, newExec := range newExecs {
		// res, err := stmt.Exec(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class, newStudent.Subject)
		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		newExec.ID = int(lastID)
		addedExecs[i] = newExec
	}
	return addedExecs, nil
}

func GetExecById(id int) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Error connecting to database ", http.StatusInternalServerError)
		return models.Exec{}, utils.ErrorHandler(err, "Error connecting to database")
	}
	defer db.Close()
	var exec models.Exec
	err = db.QueryRow("Select *from execs where id = ?", id).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.InactiveStatus, &exec.Role)
	if err == sql.ErrNoRows {
		// http.Error(w, "Student not found ", http.StatusNotFound)
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found ")
	} else if err != nil {
		// http.Error(w, "Query Error ", http.StatusInternalServerError)
		return models.Exec{}, utils.ErrorHandler(err, "Query Error ")
	}
	return exec, nil
}

func GetExecsDbOperation(execs []models.Exec, r *http.Request) ([]models.Exec, error) {

	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Error connecting to database ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Error connecting to database ")
	}
	defer db.Close()

	var args []interface{}

	query := "SELECT id, first_name, last_name, email, username, user_created_at,  inactive_status, role FROM EXECS WHERE 1=1"
	query, args = utils.AddFilters(r, query, args)
	query = utils.AddSorting(r, query)

	// if(firstName != ""){
	// 	query += " AND first_name = ?"
	// 	args = append(args,firstName)
	// }
	// if(lastName != ""){
	// 	query += " AND last_name = ?"
	// 	args = append(args, lastName)
	// }

	rows, err := db.Query(query, args...)
	if err != nil {
		// http.Error(w, "Database Query Error ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Database Query Error  ")
	}
	defer rows.Close()

	// StudentList := make([]models.Student, 0)
	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
		if err != nil {
			// http.Error(w,"Error Scanning database ",http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Error Scanning database ")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}
