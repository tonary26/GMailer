package handlers

import (
	"encoding/json"
	"gmailer/internal/database"
	"gmailer/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MailingHandler struct {
	store *database.MailingStore
}

func NewMailingHandler(store *database.MailingStore) *MailingHandler {
	return &MailingHandler{store: store}
}

func (h *MailingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input models.CreateMailingInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректные данные")
		return
	}

	mailing, err := h.store.CreateMailing(ctx, input)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Ошибка базы")
		return
	}

	respondWithJSON(w, r, http.StatusCreated, map[string]interface{}{"mailing": mailing})
}

func (h *MailingHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mailings, err := h.store.ListMailings(ctx)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректные данные")
		return
	}

	respondWithJSON(w, r, http.StatusOK, map[string]interface{}{"mailings": mailings})
}

func (h *MailingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректный ID")
		return
	}

	mailing, err := h.store.GetMailingById(ctx, id)
	if err != nil {
		respondWithError(w, r, http.StatusNotFound, "Рассылки с таким ID не найдено")
		return
	}

	respondWithJSON(w, r, http.StatusOK, map[string]interface{}{"mailing": mailing})
}

func (h *MailingHandler) Start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, "Некоректный ID")
		return
	}

	err = h.store.Start(ctx, id)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "не удалось запустить рассылку")
		return
	}

	respondWithJSON(w, r, http.StatusOK, map[string]interface{}{"status": "running"})
}
