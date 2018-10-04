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
	"postGres/models"

	_ "github.com/lib/pq"
)

// func loadTodo(id int) (*Todo, error) {
// 	row :=
// 	if err != nil {
// 			return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

var templates = template.Must(template.ParseFiles("views/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/index.html")
	// if err != nil {
	// 	fmt.Println("template error", err)
	// }
	todos, err := models.GetTodos()
	if err != nil {
		fmt.Println("query error", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = templates.ExecuteTemplate(w, "index.html", todos)
	if err != nil {
		fmt.Println("t.exec fail", err)
	}

	// data := struct{
	// 	User string,
	// 	TodosList []models.Todo
	// }{"Igor A", todos}

	// for _, todo := range todos {
	// 	fmt.Fprintf(w, "%d, %s, %d, %t\n", todo.ID, todo.Body, todo.AuthorID, todo.Done)
	// }
	//fmt.Printf("%#v", todos)
	// err = t.Execute(w, todos)
	// if err != nil {
	// 	fmt.Println("t.exec fail", err)
	// }
}

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.PATH[len("/edit/"):]

// 	t, err :=
// }

func main() {
	models.Init()
	http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/edit"), editHandler)
	http.ListenAndServe(":8000", nil)
}
