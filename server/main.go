package main

import (
	"fmt"
	"server/internal/routes"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	baseEcho := echo.New()
	baseEcho.Use(middleware.Recover())
	baseEcho.Use(middleware.RequestLogger())

	groupedEcho := baseEcho.Group("/api/v1")

	messageHanlder := routes.NewMessageRoute(groupedEcho)
	messageHanlder.RegisterRoutes()

	websocketHandler := routes.NewWebSocketRoute(groupedEcho)
	websocketHandler.RegisterRoutes()

	fmt.Println("Server started")
	if err := baseEcho.Start(":1323"); err != nil {
		baseEcho.Logger.Error("failed to start server", "error", err)
	}
}
