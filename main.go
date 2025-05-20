package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
    "github.com/Glenn444/chirpy/internal/database"
    "github.com/Glenn444/chirpy/internal/handler"
)


func Health(w http.ResponseWriter,req *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		handler.ServeHTTP(w, r)

		duration := time.Since(startTime)

		log.Printf(
			"%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			duration,
		)
	})
}


func main() {
    godotenv.Load()
    dbURL := os.Getenv("DB_URL")
	Jwt_secret := os.Getenv("SECRET")
	platform := os.Getenv("PLATFORM")
    db,err := sql.Open("postgres",dbURL)

    if err != nil{
        log.Fatal("Error Occurred in db connection")
    }
    dbQueries := database.New(db)
	cfg := &handler.ApiConfig{DB: dbQueries,Platform: platform,Secret:Jwt_secret}
	
	mux := http.NewServeMux()
	//rh := http.RedirectHandler("tobitresearchconsulting.com",307)
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/",cfg.MiddlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz",Health)
	mux.HandleFunc("GET /admin/metrics",cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset",cfg.DeleteUsers)
	// mux.HandleFunc("POST /api/validate_chirp",cfg.CreateChirps)
   
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirps)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}",cfg.GetAChirp)

	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("POST /api/login", cfg.LoginUser);
	mux.HandleFunc("POST /api/refresh", cfg.RefreshHandler);
	mux.HandleFunc("POST /api/revoke", cfg.RevokeHandler);

	loggedMux := logRequest(mux)

	server := &http.Server{
		Addr: ":8080",
		Handler: loggedMux,
	}

	log.Print("Listening...")

	if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }

}