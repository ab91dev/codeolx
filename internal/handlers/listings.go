package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/ab91dev/codeolx/internal/httpx"
	"github.com/ab91dev/codeolx/internal/middleware"
)

type listings struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price string `json:"price"`
	City string `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type ListingHandlerParams struct {
	db *sql.DB
	logger *slog.Logger
}

func NewListingHandlerParams(db *sql.DB, logger *slog.Logger) *ListingHandlerParams {
	return &ListingHandlerParams{
		db: db,
		logger: logger,
	}
}

func (lh ListingHandlerParams) List(w http.ResponseWriter, r *http.Request) {	
	// request scoped context
	ctx := r.Context()

	rows,err := lh.db.QueryContext(ctx,
		`SELECT id, title, description, price, city, created_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 10`)

	if err != nil {
		//log.Printf("query: %v",err)
		lh.logger.Error("listings query error", "err", err)
		//http.Error(w, "internal error", http.
		// StatusInternalServerError)
		httpx.Error(w,http.StatusInternalServerError,"something went wrong", httpx.CodeInternalError)
		return
	}
	defer rows.Close()

	arrListings := []listings{}

	for rows.Next() {
		var l listings
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			log.Printf("rows.scan: %v",err)
			//http.Error(w,"internal error", http.StatusInternalServerError)
			httpx.Error(w,http.StatusInternalServerError,"something went wrong", httpx.CodeInternalError)
			return
		}
		arrListings = append(arrListings,l)
	}

	if err:= rows.Err(); err != nil {
		log.Printf("rows.err; %v",err)
		//http.Error(w,"internal error", http.StatusInternalServerError)
		httpx.Error(w,http.StatusInternalServerError,"something went wrong", httpx.CodeInternalError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(arrListings)
}

func (lh ListingHandlerParams) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	id := r.PathValue("id")
	//fmt.Println("id",id)

	lh.logger.Debug("debug log","listing_id",id)
	lh.logger.Info("starting query","listing_id",id)
	lh.logger.Warn("warn log","listing_id",id)
	
	_, err := lh.db.ExecContext(ctx,
		`DELETE FROM listings WHERE id=$1`, id,
	)
	if err!=nil {
		//log.Printf("delete: %v",err)
		lh.logger.Error("delete failed","listing_id",id,"request_id",requestId,"error",err)
		//http.Error(w,"internal error", http.StatusInternalServerError)
		httpx.Error(w,http.StatusInternalServerError,"something went wrong", httpx.CodeInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
