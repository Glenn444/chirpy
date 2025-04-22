package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func Health(w http.ResponseWriter,req *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	msg := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(msg))
}

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Reset Successfuly"))
}


func main() {
	cfg := &apiConfig{}
	mux := http.NewServeMux()
	//rh := http.RedirectHandler("tobitresearchconsulting.com",307)
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/",cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("/healthz",Health)
	mux.HandleFunc("/metrics",cfg.MetricsHandler)
	mux.HandleFunc("/reset",cfg.ResetHandler)


	log.Print("Listening...")

	http.ListenAndServe(":8080", mux)

}