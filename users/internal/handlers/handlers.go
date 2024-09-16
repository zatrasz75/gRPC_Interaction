package handlers

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	rplesPb "zatrasz75/gRPC_Interaction/roles/pkg/grpc/roles"
	"zatrasz75/gRPC_Interaction/users/configs"
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

	// Вызов сервиса Roles для получения ролей
	conn, err := grpc.Dial(s.cfg.Internal.Roles_service_address, grpc.WithInsecure())
	if err != nil {
		s.l.Error("Не удалось подключиться к службе ролей", err)
		return nil, status.Error(codes.Internal, "Не удалось подключиться к службе ролей")
	}
	defer conn.Close()

	rolesClient := rplesPb.NewRolesClient(conn)
	rolesReq := &rplesPb.IdRequest{ID: payload.Id}

	rolesRes, err := rolesClient.GetRoles(c, rolesReq)
	if err != nil {
		s.l.Error("GetRoles()", err)
		return nil, status.Error(codes.Internal, "Не удалось получить роли пользователей")
	}

	prof, err := s.repo.ListOfUserProfile(id)
	if err != nil {
		s.l.Error("ListOfUserProfile()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &usersPb.ProfileUserResponse{
		Roles:      rolesRes.Roles,
		Name:       prof.Name,
		Surname:    prof.Surname,
		Patronymic: prof.Patronymic,
		Email:      prof.Email,
	}, nil
}
