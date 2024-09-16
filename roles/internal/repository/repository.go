package repository

import (
	"context"
	"fmt"
	"time"
	"zatrasz75/gRPC_Interaction/roles/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Store {
	return &Store{pg}
}

// ListOfUserRoles Возвращает список ролей пользователя по его user_id
func (s *Store) ListOfUserRoles(id int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT name FROM roles WHERE user_id = $1"

	rows, err := s.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err = rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("ошибка сканирования роли: %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %w", err)
	}

	return roles, nil
}
