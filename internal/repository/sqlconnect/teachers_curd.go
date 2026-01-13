package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/pkg/utils"
	"strconv"
	"strings"
)

func DeleteOneTeacher(id int) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Unable to connect Databases")
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Error deleting teacher ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error deleting teacher")
	}
	fmt.Println(result.RowsAffected())
	rowsEffected, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrieving delete result ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error retrieving delete result")
	}

	if rowsEffected == 0 {
		// http.Error(w, "Teacher not found ", http.StatusNotFound)
		return utils.ErrorHandler(err, "Teacher not found")
	}
	return nil
}

func DeleteTeachers(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Unable to connect Database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// http.Error(w, "Error starting transcation ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Error starting transcation")
	}
	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")

	if err != nil {
		tx.Rollback()
		// http.Error(w, "Error Preparing delete statement ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Error Preparing delete statement")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error deleting error ", http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Error Preparing delete statement")
		}
		rowAffected, err := result.RowsAffected()

		if err != nil {
			// http.Error(w, "Error retrieving deleted result", http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Error retrieving deleted result")
		}
		if rowAffected < 1 {
			tx.Rollback()
			// http.Error(w, fmt.Sprintf("Id %d does not exists ", id), http.StatusNotFound)
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Id %d does not exists ", id))
		}
		if rowAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
	}

	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Error commiting transcation", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Error commiting transcation")
	}

	if len(deletedIds) < 1 {
		// http.Error(w, "IDs do not exists ", http.StatusBadRequest)
		return nil, utils.ErrorHandler(err, "IDs do not exists")
	}
	return deletedIds, nil
}

func PatchOneTeacher(id int, updated map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Unable to connect Databases")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Teacher Not found ", http.StatusNotFound)
			return models.Teacher{}, utils.ErrorHandler(err, "Teacher Not found")
		}
		// http.Error(w, "Unable to retrieve data ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	for k, v := range updated {
		switch k {
		case "first_name":
			existingTeacher.FirstName = v.(string)
		case "last_name":
			existingTeacher.LastName = v.(string)
		case "email":
			existingTeacher.Email = v.(string)
		case "class":
			existingTeacher.Class = v.(string)
		case "subject":
			existingTeacher.Subject = v.(string)
		}
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	fmt.Println(teacherVal, "  ,, ")
	teacherType := teacherVal.Type()

	for k, v := range updated {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				fieldVal := teacherVal.Field(i)
				if fieldVal.CanSet() {
					fieldVal.Set(
						reflect.ValueOf(v).Convert(fieldVal.Type()),
					)
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject =? WHERE id = ? ", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		// http.Error(w, "Error updating teacher ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Error updating teacher")
	}
	return existingTeacher, nil
}

func PatchTeachers(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect database", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Unable to connect database")
	}
	defer db.Close()
	trx, err := db.Begin()

	if err != nil {
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error starting transaction")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			trx.Rollback()
			// http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Invalid teacher ID")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			trx.Rollback()
			// http.Error(w, "Error converting Id to int ", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Error converting Id to int")
		}

		var teacher models.Teacher
		err = trx.QueryRow(
			`SELECT id, first_name, last_name, email, class, subject 
			 FROM teachers WHERE id = ?`, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

		if err == sql.ErrNoRows {
			trx.Rollback()
			// http.Error(w, "Teacher not found", http.StatusNotFound)
			return utils.ErrorHandler(err, "Teacher not found")
		} else if err != nil {
			trx.Rollback()
			// http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Error retrieving teacher")
		}

		// Apply updates using reflection
		teacherVal := reflect.ValueOf(&teacher).Elem()
		teacherType := teacherVal.Type()
		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]

				if jsonTag == k {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							trx.Rollback()
							// http.Error(w, "Type mismatch in PATCH", http.StatusBadRequest)
							return utils.ErrorHandler(err, "Type mismatch in PATCH")
						}
					}
					break
				}
			}
		}

		_, err = trx.Exec(
			`UPDATE teachers 
			 SET first_name=?, last_name=?, email=?, class=?, subject=? 
			 WHERE id=?`,
			teacher.FirstName,
			teacher.LastName,
			teacher.Email,
			teacher.Class,
			teacher.Subject,
			teacher.ID,
		)

		if err != nil {
			trx.Rollback()
			// http.Error(w, "Error updating teacher", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Error updating teacher")
		}
	}

	if err := trx.Commit(); err != nil {
		// http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error committing transaction")
	}
	return nil
}

func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "unable to connect database ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "unable to connect database")
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT *FROM TEACHERS WHERE ID = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Unable to find teacher ", http.StatusNotFound)
			return models.Teacher{}, utils.ErrorHandler(err, "Unable to find teacher")
		}
		// http.Error(w, "unable to retrieve database ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "unable to retrieve database")

	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName,
		updatedTeacher.LastName,
		updatedTeacher.Email,
		updatedTeacher.Class,
		updatedTeacher.Subject,
		existingTeacher.ID,
	)

	if err != nil {
		// http.Error(w, "Error updating teachers ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Error updating teachers")

	}
	return updatedTeacher, nil
}


func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connect error")
	}

	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(GenerateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		// res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		values := GetStructValues(newTeacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag:", dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" { // skip the ID field if it's auto increment
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"

		}
	}
	fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func GetStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()
	values := []interface{}{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	log.Println("Values:", values)
	return values
}


func GetTeacherById(id int) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Error connecting to database ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Error connecting to database")
	}
	defer db.Close()
	var teacher models.Teacher
	err = db.QueryRow("Select *from teachers where id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Subject, &teacher.Class)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found ", http.StatusNotFound)
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher not found ")
	} else if err != nil {
		// http.Error(w, "Query Error ", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Query Error ")
	}
	return teacher, nil
}

func GetTeachersDbOperation(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {

	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Error connecting to database ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Error connecting to database ")
	}
	defer db.Close()

	var args []interface{}
	query := "SELECT *FROM TEACHERS WHERE 1=1"
	query, args = addFilters(r, query, args)
	query = addSorting(r, query)

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

	// teacherList := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			// http.Error(w,"Error Scanning database ",http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Error Scanning database ")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]

	if len(sortParams) > 0 {

		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) < 2 {
				continue
			}
			field := parts[0]
			order := parts[1]
			if !isValidField(field) || !isValidOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}

	}
	return query
}

func isValidField(field string) bool {
	validField := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"class":      true,
		"subject":    true,
		"email":      true,
	}
	return validField[field]

}

func isValidOrder(order string) bool {
	return order == "asc" || order == "desc"
}
