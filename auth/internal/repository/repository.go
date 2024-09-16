package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
	"zatrasz75/gRPC_Interaction/auth/internal/models"
	"zatrasz75/gRPC_Interaction/auth/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Store {
	return &Store{pg}
}

// UserVerification Проверяет существование пользователя
func (s *Store) UserVerification(u models.Users) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	err := s.Pool.QueryRow(ctx, query, u.Email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования пользователя: %w", err)
	}

	return exists, err
}

// CreateUser Добавляет нового пользователя
func (s *Store) CreateUser(u models.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (email, password, data) VALUES ($1, $2, $3)"
	_, err := s.Pool.Exec(ctx, query, u.Email, u.Password, u.Date)
	if err != nil {
		return fmt.Errorf("не удалось вставить пользователя: %w", err)
	}

	return nil
}

// CheckPasswordLogin Возвращает пароль пользователя пр email
func (s *Store) CheckPasswordLogin(u models.Users) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT password FROM users WHERE email=$1"
	var pass string
	err := s.Pool.QueryRow(ctx, query, u.Email).Scan(&pass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("не найдено ни одной строки")
		}
		return "", fmt.Errorf("ошибка при выполнении запроса %w", err)
	}
	return pass, nil
}

// LoginUser Проверяет авторизацию и возвращает id
func (s *Store) LoginUser(u models.Users) (models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id FROM users WHERE email=$1 AND password=$2"
	var user models.Users
	err := s.Pool.QueryRow(ctx, query, &u.Email, &u.Password).Scan(&user.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, fmt.Errorf("не найдено такого логина или пароля %w", err)
		}
		return u, fmt.Errorf("ошибка при выполнении запроса %w", err)
	}

	return user, nil
}
