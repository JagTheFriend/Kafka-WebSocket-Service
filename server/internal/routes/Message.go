// Package routes Stores the message routes
package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type MessageRoute struct {
	e *echo.Group
}

func NewMessageRoute(e *echo.Group) *MessageRoute {
	grouped := e.Group("/message")

	return &MessageRoute{
		e: grouped,
	}
}

func (h *MessageRoute) RegisterRoutes() {
	h.e.GET("/health", h.healthCheck)
}

func (h *MessageRoute) healthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "Message Route Operational")
}
