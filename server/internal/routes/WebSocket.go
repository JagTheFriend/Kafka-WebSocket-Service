// Package routes Stores the WebSocket routes
package routes

import (
	util "kafka-client/internals"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

type WebSocketRoute struct {
	e        *echo.Group
	consumer *kafka.Reader
}

func NewWebSocketRoute(e *echo.Group) *WebSocketRoute {
	grouped := e.Group("/websocket")

	return &WebSocketRoute{
		e:        grouped,
		consumer: util.NewKafkaReader("client-response", "client-response-group"),
	}
}

func (h *WebSocketRoute) RegisterRoutes() {
	h.e.GET("/health", h.healthCheck)
}

func (h *WebSocketRoute) healthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "WebSocket Route Operational")
}
