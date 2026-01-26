// Package routes Stores the WebSocket routes
package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type WebSocketRoute struct {
	e *echo.Group
}

func NewWebSocketRoute(e *echo.Group) *WebSocketRoute {
	grouped := e.Group("/websocket")

	return &WebSocketRoute{
		e: grouped,
	}
}

func (h *WebSocketRoute) RegisterRoutes() {
	h.e.GET("/health", h.healthCheck)
}

func (h *WebSocketRoute) healthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "WebSocket Route Operational")
}
