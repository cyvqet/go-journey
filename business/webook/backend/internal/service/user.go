package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrInvaildUserOrPassword = errors.New("账号/密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) error {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err == ErrUserNotFound {
		return ErrInvaildUserOrPassword
	}
	if err != nil {
		return err
	}

	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 密码错误
		// 记录日志
		return ErrInvaildUserOrPassword
	}
	return nil
}
