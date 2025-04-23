package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"encoding/json"
	"time"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	hts := cfg.fileserverHits.Load()
	msg := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,hts)

	w.Write([]byte(msg))
}

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: 0"))
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
    // Set JSON content type header
    w.Header().Set("Content-Type", "application/json")
    
    type parameters struct {
        Body string `json:"body"`
    }
    
    type errorResponse struct {
        Error string `json:"error"`
    }
    
    type successResponse struct {
        Valid bool `json:"valid"`
    }
    
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respBody := errorResponse{
            Error: "Something went wrong",
        }
        errData, err := json.Marshal(respBody)
        if err != nil {
            log.Printf("Error marshalling JSON: %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(errData)
        return
    }
    
    if len(params.Body) > 140 {
        respBody := errorResponse{
            Error: "Chirp is too long",
        }
        errData, err := json.Marshal(respBody)
        if err != nil {
            log.Printf("Error marshalling JSON: %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        w.Write(errData)
        return
    }
    
    // Valid case
    respBody := successResponse{
        Valid: true,
    }
    
    successData, err := json.Marshal(respBody)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(successData)
}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		handler.ServeHTTP(w, r)

		duration := time.Since(startTime)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			duration,
		)
	})
}

func main() {
	cfg := &apiConfig{}
	
	mux := http.NewServeMux()
	//rh := http.RedirectHandler("tobitresearchconsulting.com",307)
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/",cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz",Health)
	mux.HandleFunc("GET /admin/metrics",cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset",cfg.ResetHandler)
	mux.HandleFunc("POST /api/validate_chirp",validateHandler)

	loggedMux := logRequest(mux)

	server := &http.Server{
		Addr: ":8080",
		Handler: loggedMux,
	}

	log.Print("Listening...")

	server.ListenAndServe()

}