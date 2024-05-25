package middleware

import (
	"context"
	"credit-plus/internal/model/response"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type BasicAuthMiddleware struct {
	basicAuth string
}

func NewBasicAuthMiddleware(cfg BasicAuthConfig) *BasicAuthMiddleware {
	basicAuthVal := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.Username, cfg.Password)))
	return &BasicAuthMiddleware{
		basicAuth: "Basic " + basicAuthVal,
	}
}

type BasicAuthConfig struct {
	Username string
	Password string
}

func (b *BasicAuthMiddleware) ValidateBasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		basicAuth := r.Header.Get("Authorization")
		if basicAuth != b.basicAuth {
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resp)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resp)
			return
		}

		tokenString = strings.Split(tokenString, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("REPLACE_THIS_WITH_ENV_VAR_SECRET"), nil
		})
		if err != nil {
			slog.Error("ERROR", "error", err)
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Invalid Token"})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp)
			return
		}

		if !token.Valid {
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resp)
			return
		} else {
			claim := token.Claims.(jwt.MapClaims)
			userId := int64(claim["user"].(map[string]interface{})["Id"].(float64))
			member := claim["member"].(map[string]interface{})

			ctx := context.WithValue(r.Context(), "member", member)
			ctx = context.WithValue(ctx, "userId", userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}
