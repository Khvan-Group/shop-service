package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	customErrors "github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	dbHost := utils.GetEnv("DB_HOST")
	dbPort := utils.GetEnv("DB_PORT")
	dbName := utils.GetEnv("DB_NAME")
	dbUser := utils.GetEnv("DB_USER")
	dbPass := utils.GetEnv("DB_PASS")
	sslmode := utils.GetEnv("SSLMODE")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslmode)
	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	DB = db
	migrateSql(DB, dbName)
}

func migrateSql(db *sqlx.DB, dbName string) {
	migrationPath := "file://internal/migrations"
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		panic(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationPath, dbName, driver)
	if err != nil {
		panic(err)
	}

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}

func StartTransaction(txFunc func(*sqlx.Tx) *customErrors.CustomError) *customErrors.CustomError {
	tx, err := DB.Beginx()
	if err != nil {
		panic(err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return txFunc(tx)
}

func StartReadOnlyTransaction(txFunc func(*sqlx.Tx) *customErrors.CustomError) *customErrors.CustomError {
	tx, err := DB.BeginTxx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		panic(err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return txFunc(tx)
}
