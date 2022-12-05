package postgres

import (
	"fmt"
	"log"
	"regexp"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

const driver = "postgres"

func replaceDBName(connStr, dbName string) string {
	r := regexp.MustCompile(`dbname=([^\s]+)\s`)
	return r.ReplaceAllString(connStr, fmt.Sprintf("dbname=%s ", dbName))
}

func createTestDB(ddlConnStr string) (string, error) {
	dbName := uuid.New().String()
	ddlDB := sqlx.MustConnect(driver, ddlConnStr)
	ddlDB.MustExec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err := ddlDB.Close(); err != nil {
		return "", err
	}
	return dbName, nil
}

// create new database for testing purpose and drop the db after testing
func newTestDB(ddlConnStr, migrationDir string) (*sqlx.DB, func(), error) {
	const driver = "postgres"

	dbName, err := createTestDB(ddlConnStr)
	tearDownFn := func() {
		// drop created test database and close connection
		ddlDB, err := sqlx.Connect(driver, ddlConnStr)
		if err != nil {
			log.Fatalf("failed to connect database: %s", err.Error())
		}

		if _, err = ddlDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName)); err != nil {
			log.Fatalf("failed to drop database: %s", err.Error())
		}

		if err = ddlDB.Close(); err != nil {
			log.Fatalf("failed to close DDL database connection: %s", err.Error())
		}
	}
	if err != nil {
		return nil, tearDownFn, err
	}

	connStr := replaceDBName(ddlConnStr, dbName)
	db := sqlx.MustConnect(driver, connStr)

	tearDownFn = func() {
		// close new connection, drop created test database and close conn
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %s", err.Error())
		}
		ddlDB, err := sqlx.Connect(driver, ddlConnStr)
		if err != nil {
			log.Fatalf("failed to connect database: %s", err.Error())
		}

		if _, err = ddlDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName)); err != nil {
			log.Fatalf("failed to drop database: %s", err.Error())
		}

		if err = ddlDB.Close(); err != nil {
			log.Fatalf("failed to close DDL database connection: %s", err.Error())
		}
	}

	if err := goose.Run("up", db.DB, migrationDir); err != nil {
		return nil, tearDownFn, err
	}

	return db, tearDownFn, nil
}
