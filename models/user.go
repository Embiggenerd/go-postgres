package models

import (
	"errors"
	"fmt"
)

type User struct {
	ID        int
	Age       int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func RegisterUser(u *User) (*User, error) {
	id := 0
	sqlUser := `
		INSERT INTO users ( age, first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id; `
	err := db.QueryRow(sqlUser, u.Age, u.FirstName, u.LastName, u.Email, u.Password).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("new user id", id)
	user := new(User)
	sqlQuery := `
		SELECT * FROM users WHERE id = $1`

	row := db.QueryRow(sqlQuery, id)

	err = row.Scan(&user.ID, &user.Age, &user.FirstName,
		&user.LastName, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return user, err
}

func validatePassword(storedPassword, providedPassword string) (bool, error) {
	if storedPassword != providedPassword {
		return false, errors.New("Incorrect password")
	}
	return true, nil
}

func LoginUser(p, e string) (*User, error) {
	sqlEmailQuery := `
		SELECT * FROM users WHERE EMAIL = $1
		LIMIT 1`
	user := new(User)

	err := db.QueryRow(sqlEmailQuery, e).Scan(&user.ID, &user.Age,
		&user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		panic(errors.New("Email not found in database"))
	}
	isValid, err := validatePassword(user.Password, p)
	if err != nil && isValid != true {
		return nil, err
	}
	return user, nil
}

func GetUserFromSession(userId string) (*User, error) {
	var id int
	sqlSessionQuery := `
		SELECT * FROM sessions
		WHERE userid = $1
		LIMIT 1
		RETURNING userid`
	user := new(User)
	err := db.QueryRow(sqlSessionQuery, userId).Scan(&id)
	if err != nil {
		panic(err)
	}
	sqlUserQuery := `
		SELECT * FROM users
		WHERE id = $1`
	err = db.QueryRow(sqlUserQuery, id).Scan(&user.ID, &user.Age,
		&user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil

}
