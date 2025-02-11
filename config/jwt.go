package config

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// RSA鍵の管理
var (
	privateKey *rsa.PrivateKey
	publicKey  rsa.PublicKey
	keyID      string
	once       sync.Once
)

// 初期化
func init() {
	once.Do(func() {
		var err error
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("failed to generate RSA key: %v", err)
		}
		publicKey = privateKey.PublicKey
		keyID = generateKeyID()
	})
}

// `kid` の生成（ランダムな base64 エンコード）
func generateKeyID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func GenerateJWT(sub string, scope string) string {
	claims := jwt.MapClaims{
		"sub":   sub,
		"scope": scope,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"kid":   keyID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	signedToken, _ := token.SignedString(privateKey)
	return signedToken
}

// JWT を検証
func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return &publicKey, nil
	})
}

// JWK を取得
func GetJWKS() map[string]interface{} {
	return map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"kid": keyID,
				"alg": "RS256",
				"use": "sig",
				"n":   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
				"e":   "AQAB",
			},
		},
	}
}
