package server

import (
	"encoding/json"
	"net/http"
)

func HandleDiscovery(w http.ResponseWriter, r *http.Request) {
	discovery := map[string]interface{}{
		"issuer":                                "http://localhost:8080",
		"userinfo_endpoint":                     "http://localhost:8080/userinfo",
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discovery)
}
