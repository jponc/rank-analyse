package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type Client struct {
	jwtSecret string
}

// NewClient instantiates an Auth Client
func NewClient(jwtSecret string) (*Client, error) {
	c := &Client{
		jwtSecret: jwtSecret,
	}

	return c, nil
}

func (c *Client) GetClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.jwtSecret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
