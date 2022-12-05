package main

import (
	"log"

	"btc_billionaire/api"
	"btc_billionaire/postgres"
	"btc_billionaire/utilities"
)

func main() {
	config, err := utilities.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewStorageFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	handler := api.NewServer(db)
	handler.StartServer(config.ServerHost, config.ServerPort)
}
