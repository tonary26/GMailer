package handlers

import "github.com/go-chi/chi/v5"

func AuthRouter(h *AuthHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}

func ContactRouter(h *ContactHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/create", h.Create)
	r.Get("/list", h.List)

	return r
}