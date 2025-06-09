package handler


import (
	"net/http"
	"sync/atomic"
	"encoding/json"

	"github.com/Glenn444/chirpy/internal/database"
)

type ApiConfig struct{
	FileserverHits atomic.Int32
	DB *database.Queries
	Platform string
	Secret string
	ApiKey string
}

//middleware for metrics

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// respondWithError sends an error response back to the client
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response back to the client
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Convert the payload to JSON
	response, err := json.Marshal(payload)
	if err != nil {
		// If JSON marshaling fails, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the content type header to application/json
	w.Header().Set("Content-Type", "application/json")
	
	// Set the HTTP status code
	w.WriteHeader(code)
	
	// Write the JSON response to the client
	w.Write(response)
}