package services

import (
	"context"
	"database/sql"

	"github.com/emersion/go-imap/client"
	financetracker "github.com/iamveso/financetracker/db/sqlc"
)

type IEmailService interface {
	GetRecentMessages(ctx context.Context, count int) error
	Init(ctx context.Context, password string, imapServer string) error
	Login(ctx context.Context, password string) error
}

type IEmailRepository interface{}

type EmailConfig struct {
	Email  string
	Client *client.Client
}

type EmailRepositoryImpl struct {
	DB      *sql.DB
	Queries financetracker.Queries
}

type EmailServiceImpl struct {
	config *EmailConfig
	repo   IEmailRepository
}

func NewEmailRepository(dbConn *sql.DB) IEmailRepository {
	return &EmailRepositoryImpl{
		DB:      dbConn,
		Queries: *financetracker.New(dbConn),
	}
}

func NewEmailConfig(email string, c *client.Client) *EmailConfig {
	return &EmailConfig{
		Email:  email,
		Client: c,
	}
}

func NewEmailService(repo IEmailRepository, config *EmailConfig) IEmailService {
	return &EmailServiceImpl{
		repo:   repo,
		config: config,
	}
}

/* Repository implementation */

/* Service Implementation */
func (s *EmailServiceImpl) GetRecentMessages(ctx context.Context, count int) error {
	return nil
}

func (s *EmailServiceImpl) Init(ctx context.Context, password string, imapServer string) error {
	return nil
}

func (s *EmailServiceImpl) Login(ctx context.Context, password string) error {
	return s.config.Client.Login(s.config.Email, password)
}
