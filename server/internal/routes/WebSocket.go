// Package routes Stores the WebSocket routes
package routes

import (
	util "kafka-client/internals"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{}

type WebSocketRoute struct {
	e        *echo.Group
	consumer *kafka.Reader
}

func NewWebSocketRoute(e *echo.Group) *WebSocketRoute {
	grouped := e.Group("/websocket")

	return &WebSocketRoute{
		e:        grouped,
		consumer: util.NewKafkaReader("message", "client-response-group"),
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
	defer h.consumer.Close()

	receiverId := c.Request().Header.Get("ReceiverId")
	if receiverId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing data")
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Block until client disconnects
		select {
		case <-c.Request().Context().Done():
			return c.String(http.StatusGatewayTimeout, "Gateway timeout")
		default:
			msg, err := h.consumer.FetchMessage(c.Request().Context())
			if err != nil {
				continue
			}
			err = ws.WriteMessage(websocket.TextMessage, msg.Value)
			if err != nil {
				c.Logger().Error("failed to write WS message", "error", err)
			} else {
				h.consumer.CommitMessages(c.Request().Context(), msg)
			}
		}
	}
}
