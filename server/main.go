package main

import (
	"fmt"

	"server/internal/routes"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	baseEcho := echo.New()
	baseEcho.Use(middleware.Recover())
	baseEcho.Use(middleware.RequestLogger())

	baseEcho.Use(echoprometheus.NewMiddleware("backend-server")) // adds middleware to gather metrics
	baseEcho.GET("/metrics", echoprometheus.NewHandler())        // adds route to serve gathered metrics

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
