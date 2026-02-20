package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/hybrid-tenancy-demo/org-app/internal/config"
	"github.com/hybrid-tenancy-demo/org-app/internal/db"
)

type TenantSignupEvent struct {
	TenantID string   `json:"tenant_id"`
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

func Run(ctx context.Context, registry *db.Registry, cfg *config.Config) {
	c, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, config.AppName, sarama.NewConfig())
	if err != nil {
		log.Fatalf("kafka consumer: %v", err)
	}
	defer c.Close()

	handler := &signupHandler{registry: registry, cfg: cfg}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.Consume(ctx, []string{cfg.KafkaTopic}, handler); err != nil {
				log.Printf("consume: %v", err)
			}
		}
	}
}

type signupHandler struct {
	registry *db.Registry
	cfg      *config.Config
}

func (h *signupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *signupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *signupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var evt TenantSignupEvent
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			continue
		}
		for _, svc := range evt.Services {
			if svc == config.AppName {
				if err := db.CreateTenantDatabase(h.cfg.DBHost, h.cfg.DBPort, h.cfg.DBUser, h.cfg.DBPassword, evt.TenantID); err != nil {
					log.Printf("failed to create tenant db %s: %v", evt.TenantID, err)
					return err
				}

				log.Printf("created tenant db %s", evt.TenantID)
				if err := h.registry.Register(evt.TenantID, h.cfg.DBHost); err != nil {
					log.Printf("failed to register tenant %s: %v", evt.TenantID, err)
					return err
				} 

				log.Printf("registered tenant %s", evt.TenantID)
				log.Printf("created org-app tenant DB for %s", evt.TenantID)
				break
			}
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
