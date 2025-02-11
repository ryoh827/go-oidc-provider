package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
)

// シンプルな認可コード管理
var authCodes = struct {
	sync.Mutex
	store map[string]string
}{store: make(map[string]string)}

func generateCode() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	if clientID == "" || redirectURI == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 認可コードを発行
	code := generateCode()
	authCodes.Lock()
	authCodes.store[code] = "alice"
	authCodes.Unlock()

	http.Redirect(w, r, redirectURI+"?code="+code+"&state="+state, http.StatusFound)
}
