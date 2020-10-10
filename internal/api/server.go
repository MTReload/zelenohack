package api

import "github.com/jmoiron/sqlx"

type Server struct {
	DB *sqlx.DB
}
