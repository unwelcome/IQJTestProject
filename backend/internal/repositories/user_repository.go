package repositories

import (
	"context"
	"database/sql"

	"github.com/unwelcome/iqjtest/internal/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users(login, password_hash) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, user.Login, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*entities.UserGet, error) {
	query := `SELECT login, created_at FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	user := &entities.UserGet{ID: id}
	err := row.Scan(&user.Login, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*entities.UserGet, error) {
	query := `SELECT id, login, created_at FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.UserGet

	for rows.Next() {
		user := &entities.UserGet{}
		err = rows.Scan(&user.ID, &user.Login, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateUserLogin(ctx context.Context, id int, login string) error {
	query := `UPDATE users SET login = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, login, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, id int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
