package repository

import (
	"context"
	"fmt"
	"time"
	"zatrasz75/gRPC_Interaction/users/internal/models"
	"zatrasz75/gRPC_Interaction/users/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Store {
	return &Store{pg}
}

// ListOfUserProfile Возвращает данные пользователя по его user_id
func (s *Store) ListOfUserProfile(id int) (models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.Users
	query := "SELECT name, surname, patronymic, email FROM users WHERE user_id = $1"

	err := s.Pool.QueryRow(ctx, query, id).Scan(&u.Name, &u.Surname, &u.Patronymic, &u.Email)
	if err != nil {
		return u, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return u, nil
}
