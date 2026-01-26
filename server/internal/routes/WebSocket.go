// Package routes Stores the WebSocket routes
package routes

import (
	"common/types"
	"encoding/json"
	"fmt"
	util "kafka-client/internals"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/websocket"
)

type WebSocketRoute struct {
	e        *echo.Group
	consumer *kafka.Reader
}

func NewWebSocketRoute(e *echo.Group) *WebSocketRoute {
	grouped := e.Group("/websocket")

	return &WebSocketRoute{
		e:        grouped,
		consumer: util.NewKafkaReader("message", "message"),
	}
}

func (h *WebSocketRoute) RegisterRoutes() {
	h.e.GET("/health", h.healthCheck)
	h.e.GET("/message", h.broadCastMessage)
}

func (h *WebSocketRoute) healthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "WebSocket Route Operational")
}

func (h *WebSocketRoute) broadCastMessage(c *echo.Context) error {
	ReceiverID := c.Request().Header.Get("ReceiverId")
	if ReceiverID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing data")
	}

	websocket.Server{
		Handshake: func(cfg *websocket.Config, r *http.Request) error {
			return nil // allow all origins (dev only)
		},
		Handler: func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				msg, err := h.consumer.FetchMessage(c.Request().Context())
				if err != nil {
					c.Logger().Error("failed to fetch message", "error", err)
					continue
				}

				fmt.Printf("Key: %s\n", string(msg.Key))

				if string(msg.Key) != "message.new" {
					continue
				}

				var value types.Message
				err = json.Unmarshal(msg.Value, &value)
				if err != nil {
					continue
				}
				fmt.Printf("message: %s\n", value.ChatID)

				if value.ReceiverID != ReceiverID {
					continue
				}

				err = websocket.Message.Send(ws, msg.Value)
				if err != nil {
					c.Logger().Error("failed to write WS message", "error", err)
				} else {
					h.consumer.CommitMessages(c.Request().Context(), msg)
				}
			}
		},
	}.ServeHTTP(c.Response(), c.Request())

	return nil
}
