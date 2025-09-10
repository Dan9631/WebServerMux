package main

import (
	"log"

	"webserver.service/config"
)

func main() {
	app := config.Config{}
	error := app.Initialize(config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	if error != nil {
		log.Fatal("Failed to initialize app:", error)
	}
	app.Run()
}
