package jwt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Header struct {
	Algorithm string
	Type string
}

/*
JWT Claims (RFC 7519)
https://www.iana.org/assignments/jwt/jwt.xhtml
*/
type Claims struct {
	Issuer     string   `json:"iss"` // Issuer (e.g. Auth service)
	Expiration int64    `json:"exp"` // Expiration timestamp
	Audience   string   `json:"aud"` // Audience - service for which JWT is intended
	Subject    string   `json:"sub"` // Subject - user identity
	Name       string   `json:"name"`
	Roles      []string `json:"roles"`
}

func (c *Claims) encode() []byte {
	JSON, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal claims")

		panic(err)
	}

	b := make([]byte, base64.RawURLEncoding.EncodedLen(len(JSON)))
	base64.RawURLEncoding.Encode(b, JSON)

	return b
}

type Token struct {
	Raw       string
	Header    Header
	Claims    Claims
	Signature []byte
}

func Parse(jwt string) (*Token, error) {
	var parts []string = strings.Split(jwt, ".")

	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}

	header := new(Header)
	err := decodeBase64JSON(parts[0], header)
	if err != nil {
		return nil, errors.New("invalid headers")
	}
	if header.Algorithm != "HS256" {
		return nil, errors.New("invalid alg")
	}
	if header.Type != "JWT" {
		return nil, errors.New("invalid typ")
	}

	claims := new(Claims)
	err = decodeBase64JSON(parts[1], claims)
	if err != nil {
		return nil, errors.New("invalid claims")
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid signature")
	}

	t := new(Token)
	t.Raw = jwt
	t.Header = *header
	t.Claims = *claims
	t.Signature = signature

	return t, nil
}

func (t *Token) IsValid(secret string) bool {
	var parts []string = strings.Split(t.Raw, ".")

	expected := computeSignature(fmt.Sprintf("%s.%s", parts[0], parts[1]), secret)

	return hmac.Equal([]byte(t.Signature), expected)
}

func Sign(c Claims, secret string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := c.encode()

	message := fmt.Sprintf("%s.%s", header, payload)

	signature := computeSignature(message, secret)

	return fmt.Sprintf("%s.%s", message, base64.RawURLEncoding.EncodeToString(signature))
}
