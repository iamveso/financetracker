package main

import (
	"context"
	"log"

	"database/sql"

	"github.com/emersion/go-imap/client"
	"github.com/iamveso/financetracker/internal/handlers"
	"github.com/iamveso/financetracker/internal/services"
	"github.com/iamveso/financetracker/internal/utils"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./financetracker.db")
	if err != nil {
		log.Fatalf("failed to open db with error: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v\n", err)
	}

	ctx := context.Background()

	appPassword := utils.GetEnvOrDefault("MAIL_APP_PASSWORD", "key")
	imapServer := utils.GetEnvOrDefault("MAIL_IMAP_SERVER", "server")

	email := utils.GetEnvOrDefault("EMAIL", "email")
	c, err := client.DialTLS(imapServer, nil)
	if err != nil {
		log.Fatalf("Connecting to IMAPServer error: %v\n", err)
	}
	defer c.Logout()

	// repositories
	userRepo := services.NewUserRepository(db)

	emailRepo := services.NewEmailRepository(db)

	// services
	userService := services.NewUserService(userRepo)

	emailConfig := services.NewEmailConfig(email, c)
	emailService := services.NewEmailService(emailRepo, emailConfig)
	if err = emailService.Init(ctx, appPassword, imapServer); err != nil {
		log.Fatalf("init email service failure: %v\n", err)
	}

	handler := handlers.NewHandler(userService, emailService, emailConfig)

	go userService.RegisterUser(context.Background(), email)

	go emailService.ListenForMessages(ctx)

	if err := handler.StartServer(); err != nil {
		log.Fatalf("Server start failed with error: %v\n", err)
	}
}
