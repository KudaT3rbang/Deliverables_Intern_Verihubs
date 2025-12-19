package usecase

import (
	"context"
	"errors"
	"lendbook/internal/entity"
	"lendbook/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(ctx context.Context, email string, password string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

type userUsecase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserUsecase(repo repository.UserRepository, secret string) UserUsecase {
	return &userUsecase{
		repo:      repo,
		jwtSecret: secret,
	}
}

func (u *userUsecase) Register(ctx context.Context, email string, password string) (string, error) {
	findUser, _ := u.repo.GetByEmail(ctx, email)
	if findUser != nil {
		return "", errors.New("email address already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	newUser := entity.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := u.repo.Create(ctx, &newUser); err != nil {
		return "", err
	}

	return u.generateToken(newUser.ID)
}

func (u *userUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return u.generateToken(user.ID)
}

func (u *userUsecase) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(u.jwtSecret))
}
