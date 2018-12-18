package main

import (
	pb "github.com/aaronflower/dzone-shipping/service.user/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	key = []byte("Gakkiyui")
)

// CustomClaims is our custom metadata, which will be hashed
// and set as the second segment in or JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

// Authable defines an authorization interface.
type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

// TokenService implements the Authable interface.
type TokenService struct {
	repo Repository
}

// Decode a token string into a token object.
func (s *TokenService) Decode(token string) (*CustomClaims, error) {
	// Parse the token
	tokenType, err := jwt.ParseWithClaims(
		string(key),
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)

	// Validate the token and return the custom claims
	if claims, ok := tokenType.Claims.(*CustomClaims); ok && tokenType.Valid {
		return claims, nil
	}
	return nil, err
}

// Encode a claim into a JWT
func (s *TokenService) Encode(user *pb.User) (string, error) {
	// create teh claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "go.micro.srv.user",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token and return
	return token.SignedString(key)
}
