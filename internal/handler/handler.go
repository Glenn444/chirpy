package handler


import (
	"net/http"
	"sync/atomic"

	"github.com/Glenn444/chirpy/internal/database"
)

type ApiConfig struct{
	FileserverHits atomic.Int32
	DB *database.Queries
	Platform string
	Secret string
}

//middleware for metrics

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}