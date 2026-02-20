package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

var validServices = map[string]bool{
	"fact-app": true,
	"org-app":  true,
	"todo-app": true,
}

type TenantSignupEvent struct {
	TenantID string   `json:"tenant_id"`
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: p, topic: topic}, nil
}

func (p *Producer) PublishTenantSignup(ctx context.Context, evt TenantSignupEvent) error {
	// Filter to valid services only
	services := make([]string, 0, len(evt.Services))
	for _, s := range evt.Services {
		if validServices[s] {
			services = append(services, s)
		}
	}
	if len(services) == 0 {
		return nil
	}
	evt.Services = services
	payload, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(evt.TenantID),
		Value: sarama.ByteEncoder(payload),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
