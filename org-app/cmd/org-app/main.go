package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/hybrid-tenancy-demo/org-app/internal/config"
	"github.com/hybrid-tenancy-demo/org-app/internal/consumer"
	"github.com/hybrid-tenancy-demo/org-app/internal/db"
	"github.com/hybrid-tenancy-demo/org-app/internal/handler"
)

func main() {
	cfg := config.Load()

	registryDB, err := db.OpenTenantRegistry(cfg.TenantDBConn)
	if err != nil {
		log.Fatalf("tenant registry: %v", err)
	}
	defer registryDB.Close()
	registry := db.NewRegistry(registryDB)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go consumer.Run(ctx, registry, cfg)

	r := chi.NewRouter()
	org := handler.NewOrgHandler(registry, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)
	r.Route("/org", func(r chi.Router) {
		r.Get("/", org.List)
		r.Post("/", org.Create)
		r.Get("/{id}", org.Get)
		r.Put("/{id}", org.Update)
		r.Delete("/{id}", org.Delete)
	})

	log.Printf("org-app listening on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
