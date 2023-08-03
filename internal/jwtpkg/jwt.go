package jwtpkg

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"
	"user-management/internal/config"
	"user-management/internal/model"

	"github.com/dgrijalva/jwt-go"
)

type JwtToken struct {
	SecretKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

func NewJwtHandler(conf config.RsaPair) *JwtToken {
	privateKeyBytes, err := os.ReadFile(conf.SecretKeyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := os.ReadFile(conf.PublicKeyPath)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}

	return &JwtToken{
		SecretKey: privateKey,
		PublicKey: publicKey,
	}
}

func (j *JwtToken) MakeToken(c model.TokenClaim) (model.JwtToken, error) {

	claims := jwt.MapClaims{
		"role":     string(c.Username),
		"username": string(c.Role),
		"exp":      c.ExpirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(j.SecretKey)

	return model.JwtToken(tokenString), err
}

func (j *JwtToken) ExtractClaims(t model.JwtToken) (model.TokenClaim, error) {
	token, err := jwt.Parse(string(t), func(token *jwt.Token) (interface{}, error) {
		return j.PublicKey, nil
	})
	if err != nil {
		return model.TokenClaim{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return model.TokenClaim{}, fmt.Errorf("token has expired")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return model.TokenClaim{}, fmt.Errorf("invalid token: role claim not found")
		}
		username, ok := claims["username"].(string)
		if !ok {
			return model.TokenClaim{}, fmt.Errorf("invalid token: username claim not found")
		}

		return model.TokenClaim{
			Username: model.Username(username),
			Role:     model.Role(role),
		}, nil
	}
	return model.TokenClaim{}, fmt.Errorf("invalid token")
}
