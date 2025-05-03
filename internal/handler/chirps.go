package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Glenn444/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
    Body string `json:"body"`
    UserId uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) CreateChirps(w http.ResponseWriter, r *http.Request) {
    // Set JSON content type header
    w.Header().Set("Content-Type", "application/json")
    
   
    
    type errorResponse struct {
        Error string `json:"error"`
    }
    
    type successResponse struct {
       ID uuid.UUID `json:"id"`
       CreatedAt time.Time `json:"created_at"`
       UpdatedAt time.Time `json:"updated_at"`
       Body string `json:"body"`
       UserId uuid.UUID `json:"user_id"`
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
    data,err := validateChirp(params)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("%v\n",err)
        return
    }
    chirpParams := database.CreateChirpParams{
        Body: data,
        UserID: params.UserId,
    }

    
    createdChirp, err := cfg.DB.CreateChirp(r.Context(), chirpParams)
    if err != nil{
        fmt.Printf("Error Creating Chirp %v\n",err)
    }
    
   
    respBody := successResponse{
        ID: createdChirp.ID,
        CreatedAt: createdChirp.CreatedAt,
        UpdatedAt: createdChirp.UpdatedAt,
        Body: createdChirp.Body,
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
func validateChirp(params parameters) (string, error){
   
    
    if len(params.Body) > 140 {
        
        return "", errors.New("body length too long")
    }
    
    stxt := strings.Split(params.Body, " ")
   
    // Valid case - filter profane words
    for idx, word := range stxt {
        wordLower := strings.ToLower(word)
        if wordLower == "kerfuffle" || wordLower == "sharbert" || wordLower == "fornax" {
           stxt[idx] = "****"
        }
    }
    type successResponse struct {
        CleanedBody string `json:"cleaned_body"`
    }
    
    formattedStr := strings.Join(stxt, " ")
    respBody := successResponse{
        CleanedBody: formattedStr,
    }
    
    successData, err := json.Marshal(respBody)
    if err != nil {
        return "", errors.New("error Marshalling respBody")
    }
    
  return string(successData),nil 
}