package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"zatrasz75/gRPC_Interaction/roles/internal/models"
	"zatrasz75/gRPC_Interaction/roles/internal/payloads"
	rplesPb "zatrasz75/gRPC_Interaction/roles/pkg/grpc/roles"
)

func (s *Store) GetUserRoles(c context.Context, request *rplesPb.UserIdRequest) (response *rplesPb.UserRolesResponse, err error) {
	payload := payloads.ExtractPayload(c)
	if payload == nil {
		fmt.Println("payload = nil", payload.Id)
	}
	id, _ := strconv.Atoi(payload.Id)

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
