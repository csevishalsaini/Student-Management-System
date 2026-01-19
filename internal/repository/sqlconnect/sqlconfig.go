package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	
)

func ConnectDb() (*sql.DB, error) {
	dbname := os.Getenv("DB_NAME")
	dbUser := os.Getenv("root")
	password := os.Getenv("DB_PASSWORD")
	Host := os.Getenv("HOST")
	dbPort := os.Getenv("DB_PORT")
	fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",dbUser,password,Host,dbPort,dbname)
	fmt.Println("Connecting to DB ")
	connectingString := "root:Vikas@12345@tcp(127.0.0.1:3306)/" + dbname
	db, err := sql.Open("mysql", connectingString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to DB success ")
	return db, nil
}
