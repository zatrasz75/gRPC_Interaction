package tokens

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"zatrasz75/gRPC_Interaction/auth/internal/models"
)

func GenerateJwtToken(u models.Users, secretKey string, expiration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = u.Id
	claims["exp"] = time.Now().Add(expiration).Unix()
	// Подписываем токен с помощью секретного ключа
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
