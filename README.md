# Kafka WebSocket Service (Echo + Kafka)

This project is a **Go-based WebSocket service** built using the **Echo framework** and **Apache Kafka**.
Kafka runs inside **Docker Compose**, while the Echo application runs **locally** and communicates with Kafka for message consumption/processing.

---

## 🚀 Tech Stack

* **Go**
* **Echo (v5)** – HTTP & WebSocket server
* **Apache Kafka**
* **Docker & Docker Compose**

---

## 📐 Architecture Overview

* Kafka brokers run in **Docker Compose**
* Echo server runs **locally**
* Echo exposes **WebSocket APIs**
* Kafka consumers read messages from configured topics
* Messages can be broadcast to connected WebSocket clients

```
HTTP Client ──▶ Echo REST API ──▶ Kafka Producer ──▶ Kafka Topic
                                                       │
WebSocket Client ◀── Echo WebSocket ◀── Kafka Consumer ┘
```

---

## ⚙️ Prerequisites

Make sure you have the following installed:

* Go 1.21+
* Docker
* Docker Compose

---

## 🐳 Running Kafka (Docker Compose)

Start Kafka and Zookeeper:

```bash
cd kafka
go mod tidy
docker-compose up -d
```

Verify containers:

```bash
docker ps
```

Kafka will be available on the configured broker address (e.g. `localhost:9093`).

---

## ▶️ Running the Echo Server (Locally)

From the project root:

```bash
cd server
go mod tidy
go run main.go
```

The Echo server will start on:

```
http://localhost:1323
```

---

## 🔌 WebSocket Endpoints

### WebSocket Message Stream

```
ws://localhost:1323/api/v1/websocket/message
```

**Headers required:**

| Header     | Description       |
| ---------- | ----------------- |
| ReceiverId | Client identifier |

### Health Check

```
GET /api/v1/websocket/health
```

Response:

```
WebSocket Route Operational
```

---

### Publish Message

```
POST /api/v1/message/action
```

**Request Body (JSON):**

```json
{
  "chatId": "some-chat-id",
  "senderId": "some-user-id",
  "receiverId": "some-reciver-id",
  "content": "hello-world"
}
```

**Behavior:**

* Binds request body to `types.Message`
* Publishes message to Kafka topic: **`message`**
* Uses key: **`message.new`**

**Success Response:**

```
200 OK
Message Sent
```

---

## 📦 Kafka Configuration

* **Topics** are defined in code (example: `message`)
* **Consumer Groups** are configurable
* Kafka consumers are created using `kafka-go`

Example:

```go
util.NewKafkaReader("message", "client-response-group")
```
