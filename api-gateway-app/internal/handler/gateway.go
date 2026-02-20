package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hybrid-tenancy-demo/api-gateway-app/internal/auth"
	"github.com/hybrid-tenancy-demo/api-gateway-app/internal/proxy"
)

type Gateway struct {
	jwtKey    string
	factURL   string
	orgURL    string
	todoURL   string
}

func NewGateway(jwtKey, factURL, orgURL, todoURL string) *Gateway {
	return &Gateway{
		jwtKey:  jwtKey,
		factURL: factURL,
		orgURL:  orgURL,
		todoURL: todoURL,
	}
}

func (g *Gateway) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := auth.FromBearer(r)
		if token == "" {
			http.Error(w, "missing or invalid authorization", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ParseToken(token, g.jwtKey)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(auth.WithTenant(r.Context(), claims)))
	})
}

func (g *Gateway) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(g.AuthMiddleware)

	r.Route("/facts", func(r chi.Router) {
		r.Handle("/*", g.proxyTo(g.factURL))
	})
	r.Route("/org", func(r chi.Router) {
		r.Handle("/*", g.proxyTo(g.orgURL))
	})
	r.Route("/todos", func(r chi.Router) {
		r.Handle("/*", g.proxyTo(g.todoURL))
	})

	return r
}

func (g *Gateway) proxyTo(baseURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetTenant(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		resp, err := proxy.Forward(r, baseURL, r.URL.Path, claims.TenantID)
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		proxy.CopyResponse(w, resp)
	})
}
