package main

import (
	"github.com/magzhan/geotracker/internal/app"
	"github.com/magzhan/geotracker/pkg/config"
	"log"
)

func init() {
	err := config.Load()
	if err != nil {
		log.Fatalf("Error from config: %v", err)
	}
}
func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("Error from server: %v", err)
	}
}
