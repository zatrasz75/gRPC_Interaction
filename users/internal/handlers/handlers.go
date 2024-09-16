package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"zatrasz75/gRPC_Interaction/users/configs"
	"zatrasz75/gRPC_Interaction/users/internal/models"
	"zatrasz75/gRPC_Interaction/users/internal/payloads"
	"zatrasz75/gRPC_Interaction/users/internal/repository"
	usersPb "zatrasz75/gRPC_Interaction/users/pkg/grpc/users"
	"zatrasz75/gRPC_Interaction/users/pkg/logger"
)

type Store struct {
	repo *repository.Store
	l    logger.LoggersInterface
	cfg  *configs.Config
	usersPb.UnimplementedProfileServer
}

func New(repo *repository.Store, l logger.LoggersInterface, cfg *configs.Config) *Store {
	return &Store{
		repo: repo,
		l:    l,
		cfg:  cfg,
	}
}

func (s *Store) GetProfileUsers(c context.Context, request *usersPb.ProfileIdRequest) (response *usersPb.ProfileUserResponse, err error) {
	payload := payloads.ExtractPayload(c)
	id, _ := strconv.Atoi(payload.Id)
	fmt.Println(id)

	var u models.Users
	prof, err := s.repo.ListOfUserProfile(id)
	if err != nil {
		s.l.Error("ListOfUserRoles()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &usersPb.ProfileUserResponse{
		Roles:      u.Roles,
		Name:       prof.Name,
		Surname:    prof.Surname,
		Patronymic: prof.Patronymic,
		Email:      prof.Email,
	}, nil
}
