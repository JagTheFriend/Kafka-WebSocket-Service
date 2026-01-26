// Package routes Stores the message routes
package routes

import (
	"common/types"
	util "kafka-client/internals"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

type MessageRoute struct {
	e        *echo.Group
	producer *kafka.Writer
}

func NewMessageRoute(e *echo.Group) *MessageRoute {
	grouped := e.Group("/message")
	producer := util.NewKafkaWriter("message")

	return &MessageRoute{
		e:        grouped,
		producer: producer,
	}
}

func (h *MessageRoute) RegisterRoutes() {
	h.e.GET("/health", h.healthCheck)
	h.e.POST("/action", h.newMessage)
}

func (h *MessageRoute) healthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "Message Route Operational")
}

func (h *MessageRoute) newMessage(c *echo.Context) error {
	var user types.Message
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	ctx := c.Request().Context()
	key := "message.new"

	if err := util.WriteToKafka(ctx, h.producer, key, user); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to publish message")
	}

	return c.String(http.StatusOK, "Message Sent")
}
