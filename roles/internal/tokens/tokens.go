package tokens

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func VerifyJwtToken(tokenStr, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("не удается получить user_id из токена")
		}

		return id, nil
	} else {
		return "", errors.New("invalid token")
	}
}
