package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ab91dev/codeolx/internal/httpx"
	"github.com/ab91dev/codeolx/internal/middleware"
)

type UploadHandler struct {
	logger *slog.Logger
}

func NewUploadHandler(logger *slog.Logger) *UploadHandler {
	return &UploadHandler{
		logger: logger,
	}
}

func (uh UploadHandler) Presign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	log := uh.logger.With("request_id", requestId)

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		log.Error("no user id found in context")
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	var req PresignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode", "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	if len(req.Files) == 0 {
		log.Error("empty file")
		httpx.Error(w, http.StatusBadRequest, "files must not be empty", httpx.CodeValidationFailed)
		return
	}

	if len(req.Files) > maxImagePerListing {
		log.Error("images upload exceeded")
		httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("listing can have atmost %d images", maxImagePerListing), httpx.CodeValidationFailed)
		return
	}

	// for _, f := range req.Files {

	// }

}
