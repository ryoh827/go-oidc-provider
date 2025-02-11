package server

import (
	"encoding/json"
	"go-oidc-provider/config"
	"net/http"
	"sync"
)

// シンプルなトークン管理
var tokens = struct {
	sync.Mutex
	store map[string]string
}{store: make(map[string]string)}

func HandleToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")

	authCodes.Lock()
	user, exists := authCodes.store[code]
	if exists {
		delete(authCodes.store, code) // 使い捨て
	}
	authCodes.Unlock()

	if !exists {
		http.Error(w, "Invalid code", http.StatusBadRequest)
		return
	}

	token := config.GenerateJWT(user)

	tokens.Lock()
	tokens.store[token] = user
	tokens.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   "3600",
	})
}
