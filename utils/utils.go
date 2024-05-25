package utils

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

func Response(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	// encode the body before sending it
	if err := json.NewEncoder(w).Encode(body); err != nil {
		zap.L().Error("error encoding response body", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func GetTaskCollectionName(email string) string {
	// replace special characters with underscore
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.ToLower(re.ReplaceAllString(email, "_"))
}
