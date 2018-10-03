package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "go"
)

// GetTodos queries db for all todos objects by authorID
func GetTodos(id int) *Todo {

}

var db *sql.DB

// Init initializes our db in main
func Init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
