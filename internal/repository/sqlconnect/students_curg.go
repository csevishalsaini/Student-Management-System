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

func DeleteOneStudent(id int) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Unable to connect Databases")
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Error deleting student ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error deleting student")
	}
	fmt.Println(result.RowsAffected())
	rowsEffected, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrieving delete result ", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error retrieving delete result")
	}

	if rowsEffected == 0 {
		// http.Error(w, "Student not found ", http.StatusNotFound)
		return utils.ErrorHandler(err, "Student not found")
	}
	return nil
}

func DeleteStudents(ids []int) ([]int, error) {
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
	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")

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

func PatchOneStudent(id int, updated map[string]interface{}) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Unable to connect Databases ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "Unable to connect Databases")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Student Not found ", http.StatusNotFound)
			return models.Student{}, utils.ErrorHandler(err, "Student Not found")
		}
		// http.Error(w, "Unable to retrieve data ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	for k, v := range updated {
		switch k {
		case "first_name":
			existingStudent.FirstName = v.(string)
		case "last_name":
			existingStudent.LastName = v.(string)
		case "email":
			existingStudent.Email = v.(string)
		case "class":
			existingStudent.Class = v.(string)
		
	}

	studentVal := reflect.ValueOf(&existingStudent).Elem()
	fmt.Println(studentVal, "  ,, ")
	studentType := studentVal.Type()

	for k, v := range updated {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				fieldVal := studentVal.Field(i)
				if fieldVal.CanSet() {
					fieldVal.Set(
						reflect.ValueOf(v).Convert(fieldVal.Type()),
					)
				}
			}
		}
	}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ? ", existingStudent.FirstName, existingStudent.LastName, existingStudent.Email, existingStudent.Class, existingStudent.ID)
	if err != nil {
		// http.Error(w, "Error updating student ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "Error updating student")
	}
	return existingStudent, nil
}

func PatchStudents(updates []map[string]interface{}) error {
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
			// http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Invalid student ID")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			trx.Rollback()
			// http.Error(w, "Error converting Id to int ", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Error converting Id to int")
		}

		var student models.Student
		err = trx.QueryRow(
			`SELECT id, first_name, last_name, email, class FROM students WHERE id = ?`, id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)

		if err == sql.ErrNoRows {
			trx.Rollback()
			// http.Error(w, "Student not found", http.StatusNotFound)
			return utils.ErrorHandler(err, "Student not found")
		} else if err != nil {
			trx.Rollback()
			// http.Error(w, "Error retrieving student", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Error retrieving student")
		}

		// Apply updates using reflection
		studentVal := reflect.ValueOf(&student).Elem()
		studentType := studentVal.Type()
		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]

				if jsonTag == k {
					fieldVal := studentVal.Field(i)
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
			`UPDATE students 
			 SET first_name=?, last_name=?, email=?, class=? 
			 WHERE id=?`,
			student.FirstName,
			student.LastName,
			student.Email,
			student.Class,
			student.ID,
		)

		if err != nil {
			trx.Rollback()
			// http.Error(w, "Error updating student", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Error updating student")
		}
	}

	if err := trx.Commit(); err != nil {
		// http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error committing transaction")
	}
	return nil
}

func UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "unable to connect database ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "unable to connect database")
	}

	var existingStudent models.Student
	err = db.QueryRow("SELECT *FROM STUDENTS WHERE ID = ?", id).Scan(&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Unable to find student ", http.StatusNotFound)
			return models.Student{}, utils.ErrorHandler(err, "Unable to find student")
		}
		// http.Error(w, "unable to retrieve database ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "unable to retrieve database")

	}
	updatedStudent.ID = existingStudent.ID
	_, err = db.Exec(
		"UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		updatedStudent.FirstName,
		updatedStudent.LastName,
		updatedStudent.Email,
		updatedStudent.Class,
		existingStudent.ID,
	)

	if err != nil {
		// http.Error(w, "Error updating students ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "Error updating students")

	}
	return updatedStudent, nil
}

func AddStudentsDBHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connect error")
	}

	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", models.Student{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))
	for i, newStudent := range newStudents {
		// res, err := stmt.Exec(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class, newStudent.Subject)
		values := utils.GetStructValues(newStudent)
		for val :=range values{
			fmt.Println(val)
		}
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		newStudent.ID = int(lastID)
		addedStudents[i] = newStudent
	}
	return addedStudents, nil
}

func GetStudentById(id int) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Error connecting to database ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "Error connecting to database")
	}
	defer db.Close()
	var student models.Student
	err = db.QueryRow("Select *from students where id = ?", id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		// http.Error(w, "Student not found ", http.StatusNotFound)
		return models.Student{}, utils.ErrorHandler(err, "Student not found ")
	} else if err != nil {
		// http.Error(w, "Query Error ", http.StatusInternalServerError)
		return models.Student{}, utils.ErrorHandler(err, "Query Error ")
	}
	return student, nil
}

func GetStudentsDbOperation(students []models.Student, r *http.Request) ([]models.Student, error) {

	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Error connecting to database ", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Error connecting to database ")
	}
	defer db.Close()

	var args []interface{}
	query := "SELECT *FROM STUDENTS WHERE 1=1"
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
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			// http.Error(w,"Error Scanning database ",http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Error Scanning database ")
		}
		students = append(students, student)
	}
	return students, nil
}
