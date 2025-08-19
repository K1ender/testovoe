package main

import (
	"testovoe/internal/api"
	"testovoe/internal/config"
)

func main() {
	cfg := config.MustInit()

	app := api.New(cfg)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
