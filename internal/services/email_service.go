package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	financetracker "github.com/iamveso/financetracker/db/sqlc"
)

type IEmailService interface {
	GetRecentMessages(ctx context.Context, count int) error
	Init(ctx context.Context, password string, imapServer string) error
	ListenForMessages(ctx context.Context)
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
	if err := s.config.Client.Login(s.config.Email, password); err != nil {
		return err
	}
	return nil
}

func (s *EmailServiceImpl) ListenForMessages(ctx context.Context) {
	_, err := s.config.Client.Select("INBOX", true)
	if err != nil {
		log.Printf("Error selecting INBOX: %v\n", err)
		return
	}

	idleClient := idle.NewClient(s.config.Client)

	updates := make(chan client.Update)
	s.config.Client.Updates = updates

	for {
		stop := make(chan struct{})
		done := make(chan error, 1)

		go func() {
			done <- idleClient.Idle(stop)
		}()

		select {
		case update := <-updates:
			if _, ok := update.(*client.MailboxUpdate); ok {
				log.Println("New email arrived!")
			}

		case <-time.After(28 * time.Minute):
			close(stop)
			<-done
			log.Println("Restarting IDLE to prevent timeout...")
		}
	}
}
