// Package util Exports util methods for kafka
package util

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

var kafkaURL = "localhost:9092"

// NewKafkaReader Created Consumers
func NewKafkaReader(topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

// NewKafkaWriter Creates Producers
func NewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.Hash{},
		// BatchSize: 10,
		RequiredAcks:           kafka.RequireAll,
		Compression:            kafka.Gzip,
		BatchSize:              100,
		BatchTimeout:           10 * time.Millisecond,
		WriteTimeout:           5 * time.Second,
		ReadTimeout:            5 * time.Second,
		AllowAutoTopicCreation: true,
		Logger:                 log.Default(),
		ErrorLogger:            log.Default(),
	}
}

var (
	ErrNilWriter = errors.New("kafka writer is nil")
	ErrEmptyKey  = errors.New("kafka message key is empty")
)

// WriteToKafka sends a JSON-encoded message to Kafka safely
func WriteToKafka(
	ctx context.Context,
	writer *kafka.Writer,
	key string,
	payload any,
) error {
	if writer == nil {
		return ErrNilWriter
	}

	if key == "" {
		return ErrEmptyKey // guarantees deterministic partitioning
	}

	// Enforce timeout (protects callers)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now().UTC(),
		Headers: []kafka.Header{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "producer", Value: []byte("go-service")},
		},
	}

	if err := writer.WriteMessages(ctx, msg); err != nil {
		log.Printf(
			"kafka produce failed | topic=%s key=%s err=%v",
			writer.Topic,
			key,
			err,
		)
		return err
	}

	return nil
}
