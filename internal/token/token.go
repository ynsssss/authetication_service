package token

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken() string {
	randomBytes := make([]byte, 16)

	rand.Read(randomBytes)

	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
}
