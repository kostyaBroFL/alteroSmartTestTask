package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresConnection(
	host string,
	port int,
	username string,
	password string,
	databaseName string,
) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName))
}

func MustGetNewPostgresConnection(
	host string,
	port int,
	username string,
	password string,
	databaseName string,
) *sqlx.DB {
	if connection, err := NewPostgresConnection(
		host,
		port,
		username,
		password,
		databaseName,
	); err != nil {
		panic(fmt.Sprintf("Can not connect to postgreSql instance: %s", err.Error()))
	} else {
		return connection
	}
}
