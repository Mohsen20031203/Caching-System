package main

import (
	"chach/massager/api"
	"chach/massager/config"
	"chach/massager/db"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		return
	}
	db, err := db.NewStorege(config)
	if err != nil {
		return
	}

	server, err := api.NewServer(db, &config)
	if err != nil {
		return
	}
	server.Router.Run(":5436")
}
