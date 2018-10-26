package models

import (
	"database/sql"
	"fmt" // Package pq importing drivers for db

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "go"
	password = "postgres"
)

var db *sql.DB

// Init initializes our db in main
func Init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		host, port, user, dbname, password)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createTables()

	fmt.Println("Successfully connected!")

	// type User struct {
	// 	ID        int
	// 	Age       int
	// 	FirstName string
	// 	LastName  string
	// 	Email     string
	// }

	// sqlStatement := `SELECT * FROM users WHERE id=$1`
	// var user User

	// row := db.QueryRow(sqlStatement, 2)
	// err = row.Scan(&user.ID, &user.Age, &user.FirstName, &user.LastName, &user.Email)
	// switch err {
	// case sql.ErrNoRows:
	// 	fmt.Println("No rows were returned")
	// 	return
	// case nil:
	// 	fmt.Println(user)
	// default:
	// 	panic(err)
	// }
}

func createTables() {
	_, err := db.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			age INT,
			first_name TEXT,
			last_name TEXT,
			email TEXT UNIQUE NOT NULL,
			password TEXT
			);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Query(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			body TEXT,
			authorid INT,
			done BOOLEAN
			);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Query(`
		CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			userid INT,
			hex TEXT
			);`)
	if err != nil {
		panic(err)
	}
}
