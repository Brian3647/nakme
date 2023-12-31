package main

import (
	"os"

	"github.com/Brian3647/nakme/pkg/api"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
)

const DefaultPort = "3000"

func getPort() string {
	portEnv := os.Getenv("PORT")

	if portEnv != "" {
		return ":" + portEnv
	} else {
		return ":" + DefaultPort
	}
}

func main() {
	dotErr := godotenv.Load(); if dotErr != nil {
		panic(dotErr)
	}

	e := echo.New()
	port := getPort()
	api.AddRoutes(e)
	e.Logger.Fatal(e.Start(port))
}
