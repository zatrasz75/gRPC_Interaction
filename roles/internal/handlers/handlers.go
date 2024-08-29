package handlers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"zatrasz75/gRPC_Interaction/roles/internal/models"
	"zatrasz75/gRPC_Interaction/roles/internal/tokens"
	rplesPb "zatrasz75/gRPC_Interaction/roles/pkg/grpc/roles"
)

func (s *Store) GetUserRoles(c context.Context, request *rplesPb.UserIdRequest) (response *rplesPb.UserRolesResponse, err error) {
	md, ok := metadata.FromIncomingContext(c)
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

	idStr, err := tokens.VerifyJwtToken(bearerToken[1], s.cfg.Token.SecretKeyHere)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	id, _ := strconv.Atoi(idStr)

	var u models.Users
	u.Roles, err = s.repo.ListOfUserRoles(id)
	if err != nil {
		s.l.Error("ListOfUserRoles()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &rplesPb.UserRolesResponse{
		Roles: u.Roles,
	}, nil
}
