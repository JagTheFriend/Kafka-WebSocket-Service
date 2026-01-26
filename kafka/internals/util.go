// Package util Exports util methods for kafka
package util

import (
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

var kafkaURL = "kafka:9092"

// GetKafkaReader Created Producers
func GetKafkaReader(topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

// NewKafkaWriter Creates Consumers
func NewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		// BatchSize: 10,
	}
}
