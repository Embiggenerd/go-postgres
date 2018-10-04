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

var templates = template.Must(template.ParseFiles("views/index.html", "views/submit.html"))

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
}

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.PATH[len("/edit/"):]

// 	t, err :=
// }

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "submit.html", nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		r.ParseForm()
		fmt.Println("body:", r.Form["body"])
		todo := models.Todo{0, r.Form["body"][0], 0, false}
		fmt.Println("todo:", todo)

		_, err := models.SubmitTodo(&todo)
		if err != nil {
			fmt.Println("submitHandlerError", err)
		}
		http.Redirect(w, r, "/submit", http.StatusFound)

	}

	// todo, err := models.SubmitTodo
}

func main() {
	models.Init()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", submitHandler)
	// http.HandleFunc("/edit"), editHandler)
	http.ListenAndServe(":8000", nil)
}
