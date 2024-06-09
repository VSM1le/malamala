package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VSM1le/malamala/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// var validate *validator.Validate

//	func init() {
//		validate = validator.New()
//	}
func generateRandomstring(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateApikey() (string, error) {
	randomString, err := generateRandomstring(32)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write([]byte(randomString))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (apiCfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name" validate:"required,min=8,max=15"`
		User     string `json:"user" validate:"required,min=8,max=20"`
		Password string `json:"password" validate:"required,min=8,max=20"`
		ConPass  string `json:"conpass" validate:"required,min=8,max=20"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Validation error %v", err))
		return
	}

	if params.Password != params.ConPass {
		respondWithERROR(w, 400, fmt.Sprintf("Couldn't create user %v", err))
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)

	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateCeoUser(r.Context(), database.CreateCeoUserParams{
		ID:       uuid.New(),
		CreateAt: time.Now().UTC(),
		UpdateAt: time.Now().UTC(),
		Name:     params.Name,
		UserName: params.User,
		Password: string(bytes),
		Role:     "ceo",
	})
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Couldn't create user %v", err))
		return
	}
	respondWithJSON(w, 201, user)
}

func (apiCfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		User     string `json:"user" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user_id, err := apiCfg.DB.LoginUser(r.Context(), params.User)
	if err != nil {
		respondWithERROR(w, 400, "Username or password incorrect")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user_id.Password), []byte(params.Password))
	if err != nil {
		respondWithERROR(w, 400, "Username or password incorrect")
		return
	}
	apiKey, err := generateApikey()
	if err != nil {
		respondWithERROR(w, 500, "Error generate API key")
	}
	api, err := apiCfg.DB.GenApiKey(r.Context(), database.GenApiKeyParams{
		ID:     user_id.ID,
		ApiKey: sql.NullString{String: apiKey, Valid: true},
	})
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Can't loggin %v", err))
		return
	}

	respondWithJSON(w, 201, api)
}

func (apiCfg *apiConfig) handlerUserLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	err := apiCfg.DB.LogoutUser(r.Context(), database.LogoutUserParams{
		ID:     user.ID,
		ApiKey: sql.NullString{String: user.ApiKey.String, Valid: true},
	})
	if err != nil {
		respondWithERROR(w, 500, fmt.Sprintf("Can't logout %v", err))
	}
	respondWithJSON(w, 200, nil)
}
