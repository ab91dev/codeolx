package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ab91dev/codeolx/internal/config"
	"github.com/ab91dev/codeolx/internal/httpx"
	"github.com/ab91dev/codeolx/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	cfg config.Config
}

// Constructor
func NewAuthHandler(db *sql.DB, logger *slog.Logger, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
		cfg: cfg,
	}
}

func (ah AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {	
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	log := ah.logger.With("request_id", requestId)
	
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode","error",err)
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
		log.Error("hashing failed", "error:",err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	
	row := ah.db.QueryRowContext(ctx, 
	`INSERT INTO users (name,email,password) VALUES ($1,$2,$3)
	RETURNING id, created_at`,
	req.Name, req.Email, hash)

	var u SignUpResponse
	if err := row.Scan(&u.ID,&u.CreatedAt); err != nil {
		
		var pgErr *pgconn.PgError
		if errors.As(err,&pgErr) && pgErr.Code == "23505"{
			httpx.Error(w, http.StatusConflict, "email already taken", httpx.CodeConflict)
			return
		}

		log.Error("scanning failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := SignUpResponse{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
	}

	log.Info("new user registered","user_id",out.ID)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}

func (ah AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	log := ah.logger.With("request_id", requestId)
	
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode","error",err)
		httpx.Error(w,http.StatusBadRequest,"invalid body", httpx.CodeMalformedJSON)
		return
	}

	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	var u user
	row := ah.db.QueryRowContext(ctx,
	`SELECT id, email, password FROM users WHERE email = $1`,req.Email)

	if err := row.Scan(&u.ID, &u.Email, &u.Password); err!= nil {
		if errors.Is(err,sql.ErrNoRows){
			// dummuy hash bcrypt.CompareHashAndPassword()
			httpx.Error(w,http.StatusUnauthorized,"email or password not match",httpx.CodeUnauthenticated)
			return
		}

		log.Error("find user by email failed","err",err)
		httpx.Error(w,http.StatusInternalServerError,"something went wrong",httpx.CodeInternalError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password),[]byte(req.Password)); err != nil {
		log.Warn("password mismatch", "user_id", u.ID)
		httpx.Error(w,http.StatusUnauthorized,"email or password not match",httpx.CodeUnauthenticated)
		return
	}

	tokenTTL := 24 * time.Hour
	now := time.Now()
	claims := jwt.RegisteredClaims {
		Subject: u.ID,
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken,err := unsignedToken.SignedString([]byte(ah.cfg.JWTKey))
	if err != nil {
		log.Error("jwt signed failed","err",err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong",httpx.CodeInternalError)
		return
	}

	out := SignInResponse{
		Token: signedToken,
		ExpiresIn: int(tokenTTL.Seconds()),
	}

	log.Info("new user logged-in","user_id",u.ID)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(out)
}	
