package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// token decode
func GetTokenClaims(token string) (jwt.MapClaims, error) {

	j, _ := jwt.Parse(token, nil)
	if j == nil {
		return nil, fmt.Errorf("token parsing failed (token=%s)", token)
	}
	claims, ok := j.Claims.(jwt.MapClaims)

	if ok {
		return claims, nil
	} else {
		return nil, errors.New("get token claims failed")
	}

}

// local auth token generate
func GenerateSessionToken(secret string, second int) (string, error) {

	claims := jwt.MapClaims{}
	claims["expired_at"] = time.Now().Add(time.Second * time.Duration(second)).Unix()

	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signer.SignedString([]byte(secret))
	if err != nil {
		return token, err
	}

	return token, nil

}

// local auth token validate
func ValidateSessionToken(secret string, token string) (expired bool, err error) {

	j, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := j.Claims.(jwt.MapClaims)
	if ok && j.Valid {
		expired_at := time.Unix(int64(claims["expired_at"].(float64)), 0)
		if time.Since(expired_at).Milliseconds() > 0 {
			return true, fmt.Errorf("expired token")
		}
	}

	return false, nil

}
