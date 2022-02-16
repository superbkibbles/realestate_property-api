package main

import (
	"github.com/joho/godotenv"
	"github.com/superbkibbles/realestate_property-api/app"
)

func main() {
	godotenv.Load()
	app.StartApplication()
}
