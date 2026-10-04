package handlers

import (
	"encoding/json"
	"gmailer/internal/database"
	"gmailer/internal/models"
	"net/http"
)

type ContactHandler struct {
	store *database.ContactStore
}

func NewContactHandler(store *database.ContactStore) *ContactHandler {
	return &ContactHandler{store: store}
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.ContactCreateInput
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректные данные")
		return 
	}

	contact, err := h.store.CreateContact(ctx, input)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return 
	}

	respondWithJSON(w, r, http.StatusCreated, map[string]interface{}{"contact": contact})	
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err, contacts := h.store.ListContact(ctx)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Ошибка получения контактов")
		return 
	}

	respondWithJSON(w, r, http.StatusOK, map[string]interface{}{"contacts": contacts})
}