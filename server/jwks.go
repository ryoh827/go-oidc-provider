package server

import (
	"encoding/json"
	"go-oidc-provider/config"
	"net/http"
)

func HandleJWKS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config.GetJWKS())
}
