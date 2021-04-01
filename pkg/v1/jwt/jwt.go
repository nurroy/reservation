package jwt

import (
	"errors"
	"github.com/kenshaw/envcfg"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

// Method is the method type.
type Method int

const (
	JWT_ISSUER = `BRI AGRO`
	TYPE = `Bearer`
)

type CustomClaims struct {
	UserId     string `json:"userId,omitempty"`
	UserPass string `json:"userPass,omitempty"`
	jwt.StandardClaims
}


// func to verify fixed token
func IsFixedTokenVerified(config *envcfg.Envcfg, auth string) error {
	// Bearer token as  RFC 6750 standard
	if strings.Split(auth, " ")[0] != TYPE || strings.Split(auth, " ")[1] != config.GetKey("jwt.fixedtoken") {
		return errors.New("Invalid token")
	}

	return nil
}