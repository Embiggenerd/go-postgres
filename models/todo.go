package models 

import 

type Todo struct {
	ID int
	Body []byte
	AuthorID int
	done bool
}