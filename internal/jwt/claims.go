package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

/*
JWT Claims (RFC 7519)
https://www.iana.org/assignments/jwt/jwt.xhtml
*/
type Claims struct {
	// Registered claims
	Iss string `json:"iss"` // Issuer (e.g. Auth service)
	Exp int64  `json:"exp"` // Expiration timestamp
	Aud string `json:"aud"` // Audience - service for which JWT is intended
	Sub string `json:"sub"` // Subject - user identity

	// Public claims
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func BuildClaims() *Claims {
	return new(Claims)
}

func (c *Claims) encode() []byte {
	JSON, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal claims")

		panic(err)
	}

	return []byte(base64.RawURLEncoding.EncodeToString(JSON))

	// b := make([]byte, base64.RawURLEncoding.EncodedLen(len(JSON)))

	// base64.RawURLEncoding.Encode(b, JSON)

	// return b
}

func (c *Claims) AddIss(iss string) *Claims {
	if c.Iss != "" {
		panic("already set")
	}

	c.Iss = iss

	return c
}

func (c *Claims) AddAud(aud string) *Claims {
	if c.Aud != "" {
		panic("already set")
	}

	c.Aud = aud

	return c
}

func (c *Claims) AddSub(sub string) *Claims {
	if c.Sub != "" {
		panic("already set")
	}

	c.Sub = sub

	return c
}

func (c *Claims) AddExp(exp int64) *Claims {
	if c.Exp != 0 {
		panic("already set")
	}

	c.Exp = exp

	return c
}

func (c *Claims) AddName(name string) *Claims {
	if c.Name != "" {
		panic("already set")
	}
	c.Name = name

	return c
}

func (c *Claims) AddRoles(roles []string) *Claims {
	if len(c.Roles) != 0 {
		panic("already set")
	}
	c.Roles = roles

	return c
}
