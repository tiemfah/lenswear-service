package jwt

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type TokenClaim struct {
	UserID   string `json:"UserID"`
	UserRole string `json:"UserRole"`
	jwt.StandardClaims
}

type JWTToken interface {
	SignToken(*TokenClaim) (*string, error)
	GetClaims(tokenString string) (*TokenClaim, error)
}

type jwtToken struct {
	privatekey *rsa.PrivateKey
	publickey  *rsa.PublicKey
}

func New(signKey *rsa.PrivateKey, unsignKey *rsa.PublicKey) JWTToken {
	return &jwtToken{
		privatekey: signKey,
		publickey:  unsignKey,
	}
}

func (j *jwtToken) SignToken(tokenClaim *TokenClaim) (*string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaim)

	tokenString, err := t.SignedString(j.privatekey)
	return &tokenString, err
}

func (j *jwtToken) GetClaims(tokenString string) (*TokenClaim, error) {
	var claims TokenClaim
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return j.publickey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return &claims, nil
}
