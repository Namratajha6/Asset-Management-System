package main

import (
	"asset_management/database"
	"asset_management/servers"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		"disable"); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")

	r := server.SetupRoutes()

	log.Println("Server running on http://localhost:8006")
	if err := http.ListenAndServe(":8006", r); err != nil {
		log.Fatal(err)
	}
}
