package server

import (
	"encoding/json"
	"go-oidc-provider/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := config.ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// JWT のクレームを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// `sub` を取得
	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"sub":   sub,
		"email": "alice@example.com",
	})
}
