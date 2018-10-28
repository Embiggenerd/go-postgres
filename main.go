package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"postGres/models"
	"postGres/utils"
	"reflect"
	"strconv"

	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseFiles("views/index.html", "views/submit.html",
	"views/edit.html", "views/register.html", "views/login.html"))

type contextKey string

func authRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionHexFromCookie string
		cookie, err := r.Cookie("user-session")
		if err != nil {
			fmt.Println(err)
			err = templates.ExecuteTemplate(w, "index.html", nil)
			if err != nil {
				fmt.Println("t.exec fail", err)
			}
		} else {
			sessionHexFromCookie = cookie.Value
			user, err := models.GetUserFromSession(sessionHexFromCookie)
			if err != nil {
				fmt.Println(err)
			}

			f := func(ctx context.Context, k contextKey) {
				v := ctx.Value(k)
				if v != nil {
					fmt.Println("user value in context", v)
					return
				}
				fmt.Println("key not found:", k)
			}
			k := contextKey("user")
			ctx := context.WithValue(context.Background(), k, user)
			f(ctx, k)
			f(ctx, contextKey("color"))
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextKey("user")).(*models.User)
	if ok {
		fmt.Println("user from context works", user)
		todos, err := models.GetTodos(user.ID)
		if err != nil {
			fmt.Println("gettods fail", err)
		}
		err = templates.ExecuteTemplate(w, "index.html",
			struct{ Todos, User interface{} }{todos, user})
		if err != nil {
			fmt.Println("t.exec fail", err)
		}
	} else {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			fmt.Println("t.exec fail", err)
		}
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/"):]

	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "edit.html", id)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		r.ParseForm()
		body := r.Form["body"][0]
		fmt.Println("edit body", body)
		fmt.Println("edit id", id)

		fmt.Println("typez", reflect.TypeOf(id))

		_, err := models.EditTodo(id, body)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "submit.html", nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		user, ok := r.Context().Value(contextKey("user")).(*models.User)
		if ok {
			r.ParseForm()
			fmt.Println("body:", r.Form["body"])
			todo := models.Todo{0, r.Form["body"][0], user.ID, false}
			fmt.Println("todo:", todo)

			_, err := models.SubmitTodo(&todo)
			if err != nil {
				panic(err)
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/delete/"):]
		err := models.DeleteTodo(id)
		if err != nil {
			fmt.Println("delete error", err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		r.ParseForm()
		fmt.Println("register form:", r.Form)
		age, err := strconv.Atoi(r.Form["age"][0])
		if err != nil {
			fmt.Println(err)
		}
		user := models.User{0, age, r.Form["firstName"][0], r.Form["lastName"][0],
			r.Form["email"][0], r.Form["password"][0]}
		fmt.Println("user:", user)
		_, err = models.RegisterUser(&user)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// Validate password, if true:
//	Return user data
// 	Find old session by user id, delete
//	Create random hex string
//	Create new row in sessions table with new user id, hex
//
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		r.ParseForm()
		user, err := models.LoginUser(r.Form["password"][0], r.Form["email"][0])
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/register", http.StatusFound)
		} else {
			err = models.DeleteSession(user.ID)
			if err != nil {
				fmt.Println(err)
			}
			hex, err := utils.RandHex(10)
			if err != nil {
				fmt.Println(err)
			}
			err = models.CreateSession(hex, user.ID)
			if err != nil {
				fmt.Println(err)
			}
			cookie := &http.Cookie{
				Name:     "user-session",
				Value:    hex,
				MaxAge:   60 * 60 * 24,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "user-session",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	models.Init()
	http.HandleFunc("/", authRequired(indexHandler))
	http.HandleFunc("/submit", authRequired(submitHandler))
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/register", registerUserHandler)
	http.HandleFunc("/login", loginUserHandler)
	http.HandleFunc("/logout", logoutUserHandler)

	http.ListenAndServe(":8000", nil)
}
