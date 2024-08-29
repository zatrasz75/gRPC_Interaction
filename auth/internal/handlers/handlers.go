package handlers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"zatrasz75/gRPC_Interaction/auth/internal/models"
	"zatrasz75/gRPC_Interaction/auth/internal/tokens"
	authPb "zatrasz75/gRPC_Interaction/auth/pkg/grpc/auth"
)

func (s *Store) Register(c context.Context, request *authPb.RegisterRequest) (response *authPb.RegisterResponse, err error) {
	var u models.Users
	u.Name = request.Name
	u.Email = request.Email
	u.Password = request.Password
	u.Date = time.Now().UTC()

	existence, err := s.repo.UserVerification(u)
	if err != nil {
		s.l.Error("UserVerification()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if existence {
		s.l.Warn("пользователь уже существует %v", existence)
		return nil, status.Error(codes.InvalidArgument, "пользователь уже существует")
	}

	hashedPassword, err := s.h.HashPass(u.Password)
	if err != nil {
		s.l.Error("HashPass()", err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	u.Password = hashedPassword

	err = s.repo.CreateUser(u)
	if err != nil {
		s.l.Error("CreateUser()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authPb.RegisterResponse{
		Access: "Ok",
	}, nil
}

func (s *Store) Login(c context.Context, request *authPb.LoginRequest) (response *authPb.LoginResponse, err error) {
	var u models.Users
	u.Email = request.Email
	u.Password = request.Password

	pass, err := s.repo.CheckPasswordLogin(u)
	if err != nil {
		s.l.Error("CheckPasswordLogin()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isPassword := s.h.HashPassCheck(pass, u.Password)
	if !isPassword {
		s.l.Debug("HashPassCheck()")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else {
		u.Password = pass
	}

	user, err := s.repo.LoginUser(u)
	if err != nil {
		s.l.Error("LoginUser()", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := tokens.GenerateJwtToken(user, s.cfg.Token.SecretKeyHere, s.cfg.Token.Expiration)
	if err != nil {
		s.l.Error("GenerateJwtToken()", err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &authPb.LoginResponse{
		JwtToken: token,
	}, nil
}
