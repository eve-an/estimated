package main

import "github.com/eve-an/estimated/internal/config"

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	server, err := InitializeApp(config)
	if err != nil {
		panic(err)
	}

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
