package models

import (
	"fmt"
	"log"
)

// Todo is for storing values returns from query
type Todo struct {
	ID       int
	Body     string
	AuthorID int
	Done     bool
}

// GetTodos returns all todos in database
func GetTodos() ([]*Todo, error) {
	fmt.Println("dbz", db)

	rows, err := db.Query("SELECT * FROM todos;")
	if err != nil {
		fmt.Println("queryerror", err)
		return nil, err
	}
	fmt.Println("rowz", rows)
	defer rows.Close()

	todos := make([]*Todo, 0)

	for rows.Next() {
		todo := new(Todo)
		err := rows.Scan(&todo.ID, &todo.Body, &todo.AuthorID, &todo.Done)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
		log.Println(todo)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return todos, err
}
