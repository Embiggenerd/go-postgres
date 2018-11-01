package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"postGres/models"
	"postGres/utils"
	"strconv"

	_ "github.com/lib/pq"
)

var tmplts = template.Must(template.ParseFiles("views/index.html", "views/withoutAuth.html", "views/home.html", "views/nav.html",
	"views/head.html", "views/header.html", "views/500.html", "views/footer.html", "views/login.html", "views/editTodo.html", "views/signup.html", "views/submitTodo.html"))

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
			err = tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "withoutAuth",
				Header: "Welcome to Go Postgres Todos",
				Styles: cacheBustedCss,
				TodoId: "",
				Todos:  nil,
				User:   nil,
			})

			if err != nil {
				utils.InternalServerError(w, r)
			}

		} else {
			sessionHexFromCookie = cookie.Value

			user, err := models.GetUserFromSession(sessionHexFromCookie)
			if err != nil {
				utils.UnauthorizedUserError(w)
			}

			f := func(ctx context.Context, k contextKey) {
				v := ctx.Value(k)
				if v != nil {
					fmt.Println("user value in context", v)
					return
				}

				utils.UnauthorizedUserError(w)

			}
			k := contextKey("user")
			ctx := context.WithValue(context.Background(), k, user)
			f(ctx, k)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextKey("user")).(*models.User)
	if !ok {
		utils.InternalServerError(w, r)
	}

	todos, err := models.GetTodos(user.ID)
	if err != nil {
		utils.InternalServerError(w, r)
	}

	err = tmplts.ExecuteTemplate(w, "index.html",
		templData{
			State:  "home",
			Header: "Home",
			Styles: cacheBustedCss,
			TodoId: "",
			Todos:  todos,
			User:   user,
		})
	if err != nil {
		utils.InternalServerError(w, r)
	}

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	todoId := r.URL.Path[len("/edit/"):]

	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html",
			templData{
				State:  "editTodo",
				Header: "Edit your todo",
				Styles: cacheBustedCss,
				TodoId: todoId,
				Todos:  nil,
				User:   nil,
			})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		r.ParseForm()
		body := r.Form["body"][0]

		_, err := models.EditTodo(todoId, body)

		if err != nil {
			utils.InternalServerError(w, r)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "submitTodo",
			Header: "Submit a new todo",
			Styles: cacheBustedCss,
			TodoId: "",
			Todos:  nil,
			User:   nil,
		})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		user, ok := r.Context().Value(contextKey("user")).(*models.User)

		if !ok {
			utils.InternalServerError(w, r)
		}

		r.ParseForm()
		todo := models.Todo{
			ID:       0,
			Body:     r.Form["body"][0],
			AuthorID: user.ID,
			Done:     false,
		}

		_, err := models.SubmitTodo(&todo)

		if err != nil {
			utils.InternalServerError(w, r)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/delete/"):]
		err := models.DeleteTodo(id)
		if err != nil {
			utils.InternalServerError(w, r)
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
			utils.InternalServerError(w, r)
		}

	} else {
		r.ParseForm()
		age, err := strconv.Atoi(r.Form["age"][0])

		if err != nil {
			fmt.Println(err)
			utils.InternalServerError(w, r)
		}

		user := models.User{
			ID:        0,
			Age:       age,
			FirstName: r.Form["firstName"][0],
			LastName:  r.Form["lastName"][0],
			Email:     r.Form["email"][0],
			Password:  r.Form["password"][0]}

		_, err = models.RegisterUser(&user)

		if err != nil {
			fmt.Println(err)
			utils.InternalServerError(w, r)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func status500Handler(w http.ResponseWriter, r *http.Request) {
	tmplts.ExecuteTemplate(w, "500.html", nil)
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
			utils.InternalServerError(w, r)
		}

	} else {
		r.ParseForm()

		user, err := models.LoginUser(r.Form["password"][0], r.Form["email"][0])

		if err != nil {
			http.Redirect(w, r, "/register", http.StatusFound)

		} else {
			err = models.DeleteSession(user.ID)

			if err != nil {
				utils.InternalServerError(w, r)
			}

			hex, err := utils.RandHex(10)

			if err != nil {
				utils.InternalServerError(w, r)
			}

			err = models.CreateSession(hex, user.ID)

			if err != nil {
				utils.InternalServerError(w, r)
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

	cacheBustedCss, _ = utils.BustaCache("mainFloats.css", cacheBustedCss)

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
	http.HandleFunc("/oops", status500Handler)

	http.ListenAndServe(":8000", nil)
}
