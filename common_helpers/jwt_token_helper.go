package common_helpers

import (
	"crypto/rsa"
	"fmt"
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/golang-jwt/jwt"
)

type JwtHelper interface {
	GenerateJwtToken(claims jwt.Claims) (string, common_errors.GenericApplicationError)
	ValidateJwtToken(jwtToken string) (jwt.Claims, common_errors.GenericApplicationError)
}

type jwtHelper struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJwtHelper(privateKey string) JwtHelper {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		panic(err)
	}
	return &jwtHelper{
		privateKey: rsaPrivateKey,
		publicKey:  &rsaPrivateKey.PublicKey,
	}
}

func (helper *jwtHelper) GenerateJwtToken(claims jwt.Claims) (string, common_errors.GenericApplicationError) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	jwtToken, err := token.SignedString(helper.privateKey)
	if err != nil {
		return "", common_errors.NewGenericInternalServerError()
	}
	return jwtToken, nil
}

func (helper *jwtHelper) ValidateJwtToken(jwtToken string) (jwt.Claims, common_errors.GenericApplicationError) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return helper.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, common_errors.NewGenericUnauthorizedError()
	}
	return token.Claims, nil
}
