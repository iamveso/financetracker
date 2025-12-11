package services

import (
	"context"
	"database/sql"

	financetracker "github.com/iamveso/financetracker/db/sqlc"
	"github.com/iamveso/financetracker/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUser(ctx context.Context, email string) (UserResponse, error)
	ComparePassword(ctx context.Context, user *financetracker.User, password string) error
	HashPassword(password string) (string, error)
	RegisterUser(ctx context.Context, email string) error
}

type IUserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (financetracker.User, error)
	CreateUser(ctx context.Context, email string) error
}

type UserResponse struct {
	Email string
}

type UserRepositoryImpl struct {
	DB      *sql.DB
	Queries *financetracker.Queries
}

type UserServiceImpl struct {
	repo IUserRepository
}

func toUserResponse(user *financetracker.User) UserResponse {
	return UserResponse{
		Email: user.Email,
	}
}

func NewUserRepository(dbConn *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB:      dbConn,
		Queries: financetracker.New(dbConn),
	}
}

func NewUserService(userRepo IUserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repo: userRepo,
	}
}

/* Repository Implementation */
func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (financetracker.User, error) {
	user, err := r.Queries.GetUser(ctx, email)
	return user, err
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, email string) error {
	password, err := utils.GetEnv("PASSWORD")
	if err != nil {
		return err
	}
	return r.Queries.CreateUser(ctx, financetracker.CreateUserParams{
		Email:    email,
		Password: password,
	})
}

/* Service Implementation */
func (s *UserServiceImpl) GetUser(ctx context.Context, email string) (UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return UserResponse{}, err
	}
	return toUserResponse(&user), nil
}

func (s *UserServiceImpl) ComparePassword(ctx context.Context, user *financetracker.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (s *UserServiceImpl) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, email string) error {
	return s.repo.CreateUser(ctx, email)
}
