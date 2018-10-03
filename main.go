package main

// create http routes for adding, deleting, changing todos
// create user, todos model
// [form data --> db, db --> templates]
// incorporate sessions, login, register
// [auth middleware, data validation]
// learn testing along the way

import (
	"fmt"
	"html/template"
	"net/http"
	"postGres/models/database"

	_ "github.com/lib/pq"
)



// func loadTodo(id int) (*Todo, error) {
// 	row := 
// 	if err != nil {
// 			return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Println("template error", err)
	}
	t.Execute(w, "go todos")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.PATH[len("/edit/"):]

	t, err := 
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/edit"), editHandler)
	http.ListenAndServe(":8000", nil)
}
