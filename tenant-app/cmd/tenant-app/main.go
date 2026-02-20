package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/config"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/db"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/events"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/handler"
	"github.com/hybrid-tenancy-demo/tenant-app/internal/repo"
)

func main() {
	cfg := config.Load()

	metaDB, err := db.OpenMetaDB(cfg.MetaDBConn)
	if err != nil {
		log.Fatalf("meta db: %v", err)
	}
	defer metaDB.Close()

	producer, err := events.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("kafka producer: %v", err)
	}
	defer producer.Close()

	tenantRepo := repo.NewRepo(metaDB)
	tenantHandler := handler.NewTenantHandler(tenantRepo)
	r := chi.NewRouter()
	r.Route("/tenants", func(r chi.Router) {
		r.Post("/", handler.NewSignupHandler(tenantRepo, producer, cfg.JWTSigningKey).ServeHTTP)
		r.Get("/{id}", tenantHandler.GetTenant)
		r.Get("/", tenantHandler.GetTenants)
	})

	log.Printf("tenant-app listening on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
