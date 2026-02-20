package main

import (
	"log"
	"net/http"

	"github.com/hybrid-tenancy-demo/api-gateway-app/internal/config"
	"github.com/hybrid-tenancy-demo/api-gateway-app/internal/handler"
)

func main() {
	cfg := config.Load()
	gw := handler.NewGateway(cfg.JWTSigningKey, cfg.FactAppURL, cfg.OrgAppURL, cfg.TodoAppURL)

	log.Printf("api-gateway listening on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, gw.Routes()); err != nil {
		log.Fatal(err)
	}
}
