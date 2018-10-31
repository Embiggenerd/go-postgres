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

var tmplts = template.Must(template.ParseFiles("views/index.html", "views/withoutAuth.html", "views/home.html", "views/nav.html",
	"views/head.html", "views/header.html", "views/footer.html", "views/login.html", "views/editTodo.html", "views/signup.html", "views/submitTodo.html"))

type templData struct {
	State  string
	Header string
	Styles string
	TodoId string
	Todos  interface{}
	User   interface{}
}
type contextKey string

var cacheBustedCss string

func authRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionHexFromCookie string
		cookie, err := r.Cookie("user-session")
		if err != nil {
			fmt.Println(err)
			// err = templates.ExecuteTemplate(w, "index.html", nil)
			err = tmplts.ExecuteTemplate(w, "index.html", templData{"withoutAuth",
				"Welcome to Go Postgres Todos", cacheBustedCss, "", nil, nil,
			})

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
			fmt.Println("gettodos fail", err)
		}
		err = tmplts.ExecuteTemplate(w, "index.html",
			templData{"home", "Home", cacheBustedCss, "", todos, user})
		if err != nil {
			fmt.Println("t.exec fail", err)
		}
		// } else {
		// 	err := templates.ExecuteTemplate(w, "index.html", nil)
		// 	if err != nil {
		// 		fmt.Println("t.exec fail", err)
		// 	}
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	todoId := r.URL.Path[len("/edit/"):]

	if r.Method == "GET" {
		_, ok := r.Context().Value(contextKey("user")).(*models.User)

		if ok {
			err := tmplts.ExecuteTemplate(w, "index.html",
				templData{"editTodo", "Edit your todo", cacheBustedCss, todoId, nil, nil})

			if err != nil {
				fmt.Println("t.exec fail", err)
			}

		}
	} else {
		r.ParseForm()
		body := r.Form["body"][0]
		fmt.Println("edit body", body)
		fmt.Println("edit id", todoId)

		fmt.Println("typez", reflect.TypeOf(todoId))
		_, err := models.EditTodo(todoId, body)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{State: "submitTodo", Header: "Submit a new todo", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil})
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
		err := tmplts.ExecuteTemplate(w, "index.html", templData{
			State: "signup", Header: "Register with an email and password", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil,
		})
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
		err := tmplts.ExecuteTemplate(w, "index.html", templData{State: "login", Header: "Log in with an email and password", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil})
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

func init() {

}
func main() {
	models.Init()
	// err := utils.AmendFilename("/home/go/src/postGres/static/mainFloats.css", "hash")
	// if err != nil {
	// 	fmt.Println("renaming error", err)
	// }

	// err := utils.Visit()
	cacheBustedCss, _ = utils.BustaCache("mainFloats.css", cacheBustedCss)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println("cacheBustedCss", cacheBustedCss)
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", authRequired(indexHandler))
	http.HandleFunc("/submit", authRequired(submitHandler))
	http.HandleFunc("/edit/", authRequired(editHandler))
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/register", registerUserHandler)
	http.HandleFunc("/login", loginUserHandler)
	http.HandleFunc("/logout", logoutUserHandler)
	http.ListenAndServe(":8000", nil)
}
