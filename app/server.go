package main

import (
	"log"

	"github.com/ichi-2049/filmie-server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.RegisterRoutes(e)

	err := e.Start(":5050")
	if err != nil {
		log.Fatalln(err)
	}
}
