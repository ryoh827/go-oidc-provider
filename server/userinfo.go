package server

import (
	"encoding/json"
	"github.com/ryoh827/go-oidc-provider/config"
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	scope, _ := claims["scope"].(string)
	scopes := strings.Split(scope, ",")

	response := make(map[string]string)

	if !contains(scopes, "openid") {
		http.Error(w, "Invalid scope", http.StatusForbidden)
		return
	}

	// sub から testUser の情報を取得したことにする
	response["sub"] = sub
	if contains(scopes, "profile") {
		response["name"] = testUser["name"]
	}
	if contains(scopes, "email") {
		response["email"] = testUser["email"]
	}

	json.NewEncoder(w).Encode(response)
}

// ヘルパー関数: 配列に要素が含まれるか確認
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
