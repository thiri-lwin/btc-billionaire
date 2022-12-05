package postgres

import (
	"log"
	"os"
	"testing"

	"btc_billionaire/utilities"
)

var dbStorage *Storage

func TestMain(m *testing.M) {
	config, err := utilities.LoadConfig("../")
	if err != nil {
		log.Fatalf("failed to load config %s", err)
	}

	var tearDownFunc func()
	dbStorage, tearDownFunc, err = newTestStorageFromConfig(config)
	if err != nil {
		tearDownFunc()
		log.Fatalf("failed to connect to database %s", err)
	} else {
		defer dbStorage.GetDBConn().Close()
	}

	exitCode := m.Run()
	tearDownFunc()

	os.Exit(exitCode)
}
