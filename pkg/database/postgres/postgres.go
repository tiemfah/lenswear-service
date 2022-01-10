package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	Postgres *sql.DB
}

var dbConn = &DB{}

// ConnectPostgresSQL connection the postgres database contract
func ConnectPostgresSQL(host string, port string, username string, pass string, dbName string, sslMode bool) (*DB, error) {
	var conStr string
	if host == "" && port == "" && dbName == "" {
		return nil, errors.New("cannot estabished the connection")
	}
	if port == "APP_DATABASE_POSTGRESDB_PORT" {
		port = "5432"
	}

	if sslMode {
		conStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			host,
			port,
			username,
			pass,
			dbName)
	} else {
		conStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			username,
			pass,
			dbName)
	}

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println(err)
		return nil, err
	}
	dbConn.Postgres = db

	return dbConn, nil
}

// DisconnectPostgres disconnection the postgres database contract
func DisconnectPostgres(ctx context.Context, db *sql.DB) {
	db.Close()
	log.Println("Connect with postgres is closed")
}
