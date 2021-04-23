package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

const (
	DURATION_TOKEN         = 60 * 15     //expired (15분)
	DURATION_REFRESH_TOKEN = 60 * 15 * 2 //expired (30분)
)

func CreateToken(c *gin.Context) error {

	// acccess-token
	token, err := GenerateAccessToken(DURATION_TOKEN)
	if err != nil {
		return err
	}
	c.SetCookie("access-token", token, DURATION_TOKEN, "", "", false, false)

	// refresh-token
	token, err = GenerateRefreshToken(DURATION_REFRESH_TOKEN)
	if err != nil {
		return err
	}
	c.SetCookie("refresh-token", token, DURATION_REFRESH_TOKEN, "", "", false, false)

	return nil

}

func CleanToken(c *gin.Context) {

	c.SetCookie("access-token", "", -1, "", "", false, false)
	c.SetCookie("refresh-token", "", -1, "", "", false, false)

}

func GenerateAccessToken(duration int) (string, error) {
	return generateToken(*config.Value.SecretAccessKey, duration)
}
func GenerateRefreshToken(duration int) (string, error) {
	return generateToken(*config.Value.SecretRefreshKey, duration)
}
func ValidateAccessToken(token string) (expired bool, err error) {
	return validateToken(*config.Value.SecretAccessKey, token)
}
func ValidateRefreshToken(token string) (expired bool, err error) {
	return validateToken(*config.Value.SecretRefreshKey, token)
}

func generateToken(secret string, duration int) (string, error) {

	claims := jwt.MapClaims{}
	claims["expired_at"] = time.Now().Add(time.Millisecond * time.Duration(1000*duration)).Unix()

	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signer.SignedString([]byte(secret))
	if err != nil {
		return token, err
	}

	return token, nil

}

func validateToken(secret string, token string) (expired bool, err error) {

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
		expired := time.Unix(int64(claims["expired_at"].(float64)), 0)
		if time.Since(expired).Milliseconds() > 0 {
			return true, errors.New("expired acccess key")
		}
	}

	return false, nil

}
