package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var TokenMaxLifetimeDuration = time.Hour * 1

type Claims struct {
	UserId string `json:"user_id"`
	PairId string `json:"pair_id"`
	jwt.StandardClaims
}

func GenerateJWT(userId string, exp int64, secret string) string {
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// NOTE: basically just adds predefined values: exp=1hour
func GenerateStandardJWT(userId, secret string) string {
	exp := time.Now().Add(TokenMaxLifetimeDuration).Unix()
	return GenerateJWT(userId, exp, secret)
}

func ParseToken(token, secret string) (Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !tok.Valid {
		return Claims{}, err
	}

	return *claims, nil
}
