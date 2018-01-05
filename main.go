package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shrinkUrl/api"
	"github.com/shrinkUrl/config"
	"github.com/shrinkUrl/db"
)

func main() {
	// load configuration
	conf, _ := config.LoadConfig()

	log.Printf("Redis Host: %s", conf.Redis.Host)
	// initialize DB
	dbClient, err := db.InitConnection(&db.Config{
		Addr:     conf.Redis.Host,
		Password: "",
	})

	// exit on DB error
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}

	// initialize and start the api
	api := api.New(dbClient, conf.App.Host)
	log.Fatal(api.Start(conf.App.Port))
}
