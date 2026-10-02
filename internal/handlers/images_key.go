package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	maxImageBytes      = 5 * 1024 * 1024 // 5 mb
	maxImagePerListing = 10
	uploadPrefix       = "uploads"
	listingPrefix = "listings"
	presignTTL         = 5 * time.Minute
)

var allowedContentTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
}

func mintUploadKey(userID uuid.UUID, ext string) string {
	return fmt.Sprintf("%s/%s/%s.%s", uploadPrefix, userID, uuid.New(), ext)
}

func mintFinalObjectKey(listingID uuid.UUID, imageID uuid.UUID, ext string) string {
	return fmt.Sprintf("%s/%s/%s.%s", listingPrefix, listingID, imageID, ext)
}

func parseUploadKey(key string, userID uuid.UUID) (uuid.UUID,string,error) {
	str,ok := strings.CutPrefix(key, fmt.Sprintf("%s/%s/",uploadPrefix,userID))
	if !ok {
		return uuid.Nil,"", errors.New("object_key does not belong to this user")
	}

	if strings.Contains(str, "/") {
		return uuid.Nil,"", errors.New("object_key has unexpected shape")
	}

	base,ext,ok := strings.Cut(str, ".")
	if !ok || !allowedExt(ext)  {
		return uuid.Nil,"", errors.New("object_key has unexpected key")
	}

	id,err := uuid.Parse(base)
	if err != nil {
		return uuid.Nil,"", errors.New("object_key is not a well-formed upload key")
	}

	return id, ext, nil
}

func allowedExt(ext string) bool {
	for _, e := range allowedContentTypes {
		if e == ext {
			return true
		}
	}
	return false
}
