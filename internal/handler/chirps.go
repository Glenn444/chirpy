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
        UserId: createdChirp.UserID,
    }
    
    successData, err := json.Marshal(respBody)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    w.Write(successData)
}

func (cfg *ApiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request)  {
     // Set JSON content type header
     w.Header().Set("Content-Type", "application/json")
    
     
     type successResponse struct {
        ID uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body string `json:"body"`
        UserId uuid.UUID `json:"user_id"`
     }
   
     var resp []successResponse

     allChirps, err := cfg.DB.GetAllChirps(r.Context())
     if err != nil{
         fmt.Printf("Error Creating Chirp %v\n",err)
     }
     
    for _,chirp := range allChirps{
        newResp := successResponse{
            ID: chirp.ID,
            CreatedAt: chirp.CreatedAt,
            UpdatedAt: chirp.UpdatedAt,
            Body: chirp.Body,
            UserId: chirp.UserID,
        }
        resp = append(resp,newResp)
    }
     
     
     successData, err := json.Marshal(resp)
     if err != nil {
         log.Printf("Error marshalling JSON: %s", err)
         w.WriteHeader(http.StatusInternalServerError)
         return
     }
     
     w.WriteHeader(http.StatusOK)
     w.Write(successData)
}

func (cfg *ApiConfig) GetAChirp(w http.ResponseWriter, r *http.Request)  {
    // Set JSON content type header
    w.Header().Set("Content-Type", "application/json")
   
    paramId, err := uuid.Parse(r.PathValue("chirpID"))
    if err != nil{
        fmt.Printf("Error Parsing uuid\n")
        return
    }
  
    type successResponse struct {
       ID uuid.UUID `json:"id"`
       CreatedAt time.Time `json:"created_at"`
       UpdatedAt time.Time `json:"updated_at"`
       Body string `json:"body"`
       UserId uuid.UUID `json:"user_id"`
    }
  
    aChirp, err := cfg.DB.GetChirp(r.Context(),paramId)
    if err != nil{
        fmt.Printf("Error Getting Chirp %v\n",err)
    }
    if aChirp.ID == uuid.Nil{
        w.WriteHeader(http.StatusNotFound)
        return
    }

    resp := successResponse{
        ID: aChirp.ID,
        CreatedAt: aChirp.CreatedAt,
        UpdatedAt: aChirp.UpdatedAt,
        Body: aChirp.Body,
        UserId: aChirp.UserID,
    }
    
   
    
    successData, err := json.Marshal(resp)
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

    
    formattedStr := strings.Join(stxt, " ")
   
    
  return formattedStr,nil 
}