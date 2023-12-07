package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Sign(c Claims, secret string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := c.encode()

	message := fmt.Sprintf("%s.%s", header, payload)

	signature := computeSignature(message, secret)

	return fmt.Sprintf("%s.%s", message, base64.RawURLEncoding.EncodeToString(signature))
}

func computeSignature(message, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	return mac.Sum(nil)
}
