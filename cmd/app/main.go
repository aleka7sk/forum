package main

import (
	"forum/config"
	"forum/server"
	"log"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Config initialize error: %v", err)
	}
	app := server.NewApp(config)

	if err := app.Run(config); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
