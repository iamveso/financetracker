package handlers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iamveso/financetracker/internal/services"
	"github.com/iamveso/financetracker/internal/utils"
)

var tpl = template.Must(template.ParseFiles("resources/html/index.html"))

type Handler struct {
	userService  services.IUserService
	emailService services.IEmailService
	emailConfig  *services.EmailConfig
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewHandler(userService services.IUserService, emailService services.IEmailService, emailConfig *services.EmailConfig) *Handler {
	return &Handler{
		userService:  userService,
		emailService: emailService,
		emailConfig:  emailConfig,
	}
}

func (h *Handler) StartServer() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/css"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	})

	// Authentication routes
	r.Post("/login", h.UserLogin)
	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dashboard Page"))
	})

	return http.ListenAndServe(":8080", r)
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	expectedPassword := utils.GetEnvOrDefault("PASSWORD", "")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email != h.emailConfig.Email || password != expectedPassword {
		w.Write([]byte("<div>Invalid email or password</div>"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// redirect to dashboard
	w.Header().Set("HX-Redirect", "/dashboard")
}
