package models

import "time"

type User struct {
	ID 			    int64			`json:"id" db:"id"`
	Email		    string		`json:"email" db:"email"`
	PasswordHash    string		`json:"password_hash"`
	CreatedAt	    time.Time	`json:"created_at" db:"created_at"`
}

type Register struct {
	Email		string		`json:"email" validate:"required,email"`
	Password    string		`json:"password" validate:"required"`
}

type Login struct {
	Email		string		`json:"email" validate:"required,email"`
	Password    string		`json:"password" validate:"required"`
}

type UserCreateInput struct {
	Email 		string	
	Password    string		
}