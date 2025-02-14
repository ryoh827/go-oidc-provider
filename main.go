package main

import (
	"github.com/ryoh827/go-oidc-provider/server"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/authorize", server.HandleAuthorize)
	http.HandleFunc("/token", server.HandleToken)
	http.HandleFunc("/userinfo", server.HandleUserInfo)
	http.HandleFunc("/.well-known/openid-configuration", server.HandleDiscovery)
	http.HandleFunc("/.well-known/jwks.json", server.HandleJWKS)

	log.Println("OIDC Provider is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
