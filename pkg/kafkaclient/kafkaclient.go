package kafkaclient

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// Client wraps Kafka producer for event publishing
type Client struct {
	producer sarama.SyncProducer
	brokers  []string
}

// NewClient creates a new Kafka client
func NewClient(brokers []string) (*Client, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Client{
		producer: producer,
		brokers:  brokers,
	}, nil
}

// PublishEvent publishes an event to a Kafka topic
func (c *Client) PublishEvent(topic string, event interface{}) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := c.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Event published to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

// Close closes the Kafka producer
func (c *Client) Close() error {
	return c.producer.Close()
}

// Event types for different services
type UserCreatedEvent struct {
	UserID    uint64 `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserUpdatedEvent struct {
	UserID    uint64 `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UpdatedAt string `json:"updated_at"`
}

type ProductCreatedEvent struct {
	ProductID uint64  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

type TaskCreatedEvent struct {
	TaskID    uint64 `json:"task_id"`
	UserID    uint64 `json:"user_id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}


