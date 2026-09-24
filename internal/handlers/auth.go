package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ab91dev/codeolx/internal/httpx"
	"github.com/ab91dev/codeolx/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// Auth Handler
type AuthHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

// Constructor
func NewAuthHandler(db *sql.DB, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
	}
}

func (ah AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {	
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("failed to decode","request_id", requestId,"error",err)
		httpx.Error(w,http.StatusBadRequest,"invalid body", httpx.CodeMalformedJSON)
		return
	}

	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	// cost = how slow you want your hashing function - 2^cost rounds.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		ah.logger.Error("hashing failed", "error:",err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	
	row := ah.db.QueryRowContext(ctx, 
	`INSERT INTO users (name,email,password) VALUES ($1,$2,$3)
	RETURNING id, created_at`,
	req.Name, req.Email, hash)

	var u user
	if err := row.Scan(&u.ID,&u.CreatedAt); err != nil {
		ah.logger.Error("scanning failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

}

