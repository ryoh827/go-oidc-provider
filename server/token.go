package server

import (
	"encoding/json"
	"github.com/ryoh827/go-oidc-provider/config"
	"net/http"
	"sync"
)

// トークン管理
var tokens = struct {
	sync.Mutex
	store map[string]string
}{store: make(map[string]string)}

// トークンエンドポイント
func HandleToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")

	authCodes.Lock()
	data, exists := authCodes.store[code]
	if exists {
		delete(authCodes.store, code) // 使い捨て
	}
	authCodes.Unlock()

	if !exists {
		http.Error(w, "Invalid code", http.StatusBadRequest)
		return
	}

	token := config.GenerateJWT(data.sub, data.scope)

	tokens.Lock()
	tokens.store[token] = data.sub
	tokens.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   "3600",
		"scope":        data.scope, // スコープを含める
	})
}
