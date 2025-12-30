package postgres

import (
	"context"
	"errors"
	"lendbook/internal/db"
	"lendbook/internal/entity"
	"lendbook/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	queries *generated.Queries
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{
		queries: generated.New(db),
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	params := generated.CreateUserParams{
		Email:    user.Email,
		Password: user.Password,
	}

	createdUser, err := u.queries.CreateUser(ctx, params)
	if err != nil {
		return errors.New("failed to create user")
	}

	user.ID = int(createdUser.ID)
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	userModel, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("failed to fetch user")
		}
		return nil, errors.New("failed to fetch user")
	}

	return &entity.User{
		ID:       int(userModel.ID),
		Email:    userModel.Email,
		Password: userModel.Password,
	}, nil
}

func (u *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	userModel, err := u.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("failed to fetch user")
		}
		return nil, errors.New("failed to fetch user")
	}

	return &entity.User{
		ID:       int(userModel.ID),
		Email:    userModel.Email,
		Password: userModel.Password,
	}, nil
}
