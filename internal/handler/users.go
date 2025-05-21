package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Glenn444/chirpy/internal/auth"
	"github.com/Glenn444/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
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
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Fatalf("Error creating user %v\n", err)
	}

	newUser := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := respBody{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
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

func (cfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email              string        `json:"email"`
		Password           string        `json:"password"`
		Expires_in_seconds time.Duration `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if params.Expires_in_seconds == time.Duration(0) || params.Expires_in_seconds > 3600 {

		params.Expires_in_seconds = 3600
	}

	type respBody struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email or password"))
		fmt.Printf("User not found: %v", err)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email or password"))
		fmt.Printf("User not found: %v", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.Secret, params.Expires_in_seconds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating jwt token"))
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Printf("Error generating refreshToken: %v\n", err)
		return
	}
	expiresAt := time.Now().AddDate(0, 0, 60)
	newRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	_, err = cfg.DB.CreateRefreshToken(r.Context(), newRefreshTokenParams)
	if err != nil {
		fmt.Printf("Error in saving refreshtoken in db: %v\n", err)
	}
	resp := respBody{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
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
	if cfg.Platform != "dev" {
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

func (cfg *ApiConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization Header required")
		return
	}
	//Check header format
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
		return
	}
	refreshToken := headerParts[1]
	fmt.Printf("Token in refreshHandler: %v \n", refreshToken)
	//Validate refresh Token
	userId, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh Token")
		return
	}

	//Generate new access token
	accessToken, err := auth.MakeJWT(userId, cfg.Secret, time.Duration(3600))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}
	type RefreshResponse struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, RefreshResponse{
		Token: accessToken,
	})
}

// RevokeHandler handles the POST /api/revoke endpoint
func (cfg *ApiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	// Check header format
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}
	refreshToken := headerParts[1]

	// Revoke token
	err := cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
