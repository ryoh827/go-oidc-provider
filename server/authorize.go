package server

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
)

// 認可コードとスコープの管理
var authCodes = struct {
	sync.Mutex
	store map[string]struct {
		sub   string
		scope string
	}
}{store: make(map[string]struct {
	sub   string
	scope string
})}

// UserInfo for test(Alice)
var testUser = map[string]string{
	"sub":   "1234567890",
	"name":  "Alice",
	"email": "alice@example.com",
}

// ランダムな認可コードを生成
func generateCode() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// 認可エンドポイント
func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")
	scope := r.URL.Query().Get("scope") // 追加

	if clientID == "" || redirectURI == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 認可コードを発行し、スコープと紐付けて保存
	code := generateCode()
	authCodes.Lock()
	authCodes.store[code] = struct {
		sub   string
		scope string
	}{testUser["sub"], scope}
	authCodes.Unlock()

	log.Println(code)

	// 認可コードをリダイレクトで返す
	http.Redirect(w, r, redirectURI+"?code="+code+"&state="+state, http.StatusFound)
}
