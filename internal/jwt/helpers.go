package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

func computeSignature(message, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	return mac.Sum(nil)
}

func decodeBase64JSON[TDst any](s string, dst TDst) error {
	decoded, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	err = json.Unmarshal(decoded, dst)
	if err != nil {
		return err
	}

	return nil
}
