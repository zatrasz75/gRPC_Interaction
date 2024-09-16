package tokens

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// AuthInterceptor проверяет JWT токен, извлекает данные и добавляет их в контекст.
func AuthInterceptor(secretKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !strings.Contains(info.FullMethod, "/users.Profile") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "метаданные не предоставляются")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "токен авторизации не предоставлен")
		}

		bearerToken := strings.Split(authHeader[0], " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "недопустимый формат токена")
		}

		id, err := verifyJwtToken(bearerToken[1], secretKey)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Неавторизованный")
		}

		pairs := []string{
			"id", id,
		}
		newMD := metadata.Pairs(pairs...)
		newCtx := metadata.NewIncomingContext(ctx, newMD)

		return handler(newCtx, req)
	}
}

func verifyJwtToken(tokenStr, secretKey string) (string, error) {
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
