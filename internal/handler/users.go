package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Glenn444/chirpy/internal/auth"
	"github.com/Glenn444/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	type respBody struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	hashedPassword,err := auth.HashPassword(params.Password);
	if err != nil{
		log.Fatalf("Error creating user %v\n",err)
	}
	
	newUser := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	}


	user, err := cfg.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	resp := respBody{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	successData, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error Marshalling %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(successData)
}

func (cfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request)  {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	type respBody struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email);
	if err != nil{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email or password"))
		fmt.Printf("User not found: %v",err);
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword,params.Password);
	if err != nil{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email or password"))
		fmt.Printf("User not found: %v",err);
		return
	}
	
	resp := respBody{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	successData, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error Marshalling %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(successData)
}




func (cfg *ApiConfig) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	

	err := cfg.DB.DeleteUsers(r.Context())
	if err != nil {
		fmt.Printf("Error Deleting users: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(successData)
}