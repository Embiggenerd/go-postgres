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

// SubmitTodo inserts values into todo table, querys by returned id, returns added todo
func SubmitTodo(t *Todo) (*Todo, error) {
	id := 0
	sqlInsert := `
		INSERT INTO todos ( body, authorId, done)
		VALUES ($1, $2, $3)
		RETURNING id`
	err := db.QueryRow(sqlInsert, t.Body, t.AuthorID, t.Done).Scan(&id)
	// err := db.QueryRow(sqlInsert, "train elephants", 123, false).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("new id", id)
	todo := new(Todo)
	sqlQuery := `SELECT * FROM todos WHERE id = $1`
	row := db.QueryRow(sqlQuery, id)

	err = row.Scan(&todo.ID, &todo.Body, &todo.AuthorID, &todo.Done)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// GetTodos returns all todos in database
func GetTodos() ([]*Todo, error) {

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
