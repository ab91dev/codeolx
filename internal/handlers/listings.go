package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
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
}

func NewListingHandlerParams(db *sql.DB) *ListingHandlerParams {
	return &ListingHandlerParams{
		db: db,
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
		log.Printf("query: %v",err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	arrListings := []listings{}

	for rows.Next() {
		var l listings
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			log.Printf("rows.scan: %v",err)
			http.Error(w,"internal error", http.StatusInternalServerError)
			return
		}
		arrListings = append(arrListings,l)
	}

	if err:= rows.Err(); err != nil {
		log.Printf("rows.err; %v",err)
		http.Error(w,"internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(arrListings)
}

func (lh ListingHandlerParams) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	//fmt.Println("id",id)
	
	_, err := lh.db.ExecContext(ctx,
		`DELETE FROM listings WHERE id=$1`, id,
	)
	if err!=nil {
		log.Printf("delete: %v",err)
		http.Error(w,"internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
