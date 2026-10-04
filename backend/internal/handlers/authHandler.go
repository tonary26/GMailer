package handlers

import (
	"encoding/json"
	"gmailer/internal/database"
	"gmailer/internal/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store *database.AuthStore
}

func NewAuthHandler(store *database.AuthStore) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.Register
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректные данные")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Ошибка хешерования пароля")
		return
	}

	createInput := models.UserCreateInput{
		Email:    input.Email,
		Password: string(hash),
	}

	user, err := h.store.Register(ctx, createInput)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.Login
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректные данные")
		return
	}

	user, err := h.store.GetByEmail(ctx, input.Email)
	if err != nil || user == nil {
		respondWithError(w, r, http.StatusBadRequest, "Неправильный логин или пароль")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Неправильный логин или пароль")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Ошибка генерации токена")
		return
	}

	respondWithJSON(w, r, http.StatusOK, map[string]interface{}{"token": tokenString, "user": user})
}
