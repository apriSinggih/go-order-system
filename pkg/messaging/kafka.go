package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

const (
	TopicOrderPaid = "order.paid"
)

func NewKafkaProducer(brokerAddr string) (*KafkaProducer, error) {
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke Kafka Broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan controller Kafka: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke Kafka controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             TopicOrderPaid,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	_ = controllerConn.CreateTopics(topicConfigs...)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddr),
		Topic:    TopicOrderPaid,
		Balancer: &kafka.LeastBytes{},
	}

	log.Println("Berhasil terhubung ke Kafka Broker!")
	return &KafkaProducer{writer: writer}, nil
}

func (k *KafkaProducer) PublishEvent(ctx context.Context, key string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal event payload: %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = k.writer.WriteMessages(ctxTimeout, kafka.Message{
		Key:   []byte(key),
		Value: body,
		Time:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("gagal publish message ke Kafka: %w", err)
	}

	log.Printf("[Kafka] Event dipublish ke topic '%s' (Key: %s): %s", k.writer.Topic, key, string(body))
	return nil
}

func (k *KafkaProducer) Close() {
	if k.writer != nil {
		k.writer.Close()
	}
}