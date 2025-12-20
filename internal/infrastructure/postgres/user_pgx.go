package postgres

import (
	"context"
	"errors"
	"lendbook/internal/entity"
	"lendbook/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`

	err := u.db.QueryRow(ctx, query, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`

	var user entity.User
	err := u.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, errors.New("failed to fetch user")
	}

	return &user, nil
}

func (u *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `SELECT id, email, password FROM users WHERE id = $1`

	var user entity.User
	err := u.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, errors.New("failed to fetch user")
	}

	return &user, nil
}
