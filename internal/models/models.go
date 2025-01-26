package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	User   UserModel
	Tokens TokenModel
}

func New(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{
			DB: db,
		},
		User: UserModel{
			DB: db,
		},
		Tokens: TokenModel{
			DB: db,
		},
	}
}
