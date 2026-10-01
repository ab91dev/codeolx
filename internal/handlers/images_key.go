package handlers

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	maxImageBytes      = 5 * 1024 * 1024 // 5 mb
	maxImagePerListing = 10
	uploadPrefix       = "uploads"
)

var allowedContentTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
}

func mintUploadKey(userID uuid.UUID, ext string) string {
	return fmt.Sprintf("%s/%s/%s.%s", uploadPrefix, userID, uuid.New(), ext)
}
