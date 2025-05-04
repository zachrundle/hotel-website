package dbrepo

import (
	"database/sql"

	"github.com/zachrundle/hotel-website/internal/config"
	"github.com/zachrundle/hotel-website/internal/repository"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB *sql.DB
}


func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: a,
		DB: conn,
	}
}