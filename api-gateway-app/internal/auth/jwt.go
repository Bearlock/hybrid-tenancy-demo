package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const TenantContextKey contextKey = "tenant"

type TokenClaims struct {
	jwt.RegisteredClaims
	TenantID string   `json:"tenant_id"`
	Services []string `json:"services,omitempty"`
}

func FromBearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(h[7:])
}

func ParseToken(raw string, signingKey string) (*TokenClaims, error) {
	t, err := jwt.ParseWithClaims(raw, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*TokenClaims)
	if !ok || !t.Valid {
		return nil, err
	}
	return claims, nil
}

func WithTenant(ctx context.Context, claims *TokenClaims) context.Context {
	return context.WithValue(ctx, TenantContextKey, claims)
}

func GetTenant(ctx context.Context) *TokenClaims {
	v := ctx.Value(TenantContextKey)
	if v == nil {
		return nil
	}
	return v.(*TokenClaims)
}
