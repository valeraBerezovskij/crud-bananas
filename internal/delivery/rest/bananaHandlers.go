package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"context"
	"strconv"
	"valerii/crudbananas/internal/domain"
)

func (h *Handler) routeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/items/"):]
	if path == "" {
		h.redirect(w, r)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		logError("routeHandler", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "id", id))
	h.redirectWithID(w, r)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	case http.MethodGet:
		h.getAll(w, r)
	}
}

func (h *Handler) redirectWithID(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		http.Error(w, "ID not found in context", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Printf("Fetching item with ID %d\n", id)
		h.getByID(w, r)
	case http.MethodPut:
		fmt.Printf("Updating item with ID %d\n", id)
		h.update(w, r)
	case http.MethodDelete:	
		fmt.Printf("Deleting item with ID %d\n", id)
		h.delete(w, r)
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var banana domain.Banana

	if err := json.NewDecoder(r.Body).Decode(&banana); err != nil {
		logError("create", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createdBanana, err := h.bananaService.Create(banana)
	if err != nil {
		logError("create", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdBanana); err != nil {
		logError("create", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	bananas, err := h.bananaService.GetAll()
	if err != nil {
		logError("getAll", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(bananas); err != nil {
		logError("getAll", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		logError("delete", fmt.Errorf("ID not found in context"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	banana, err := h.bananaService.GetById(id)
	if err != nil {
		logError("getByID", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(banana); err != nil {
		logError("getByID", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		logError("update", fmt.Errorf("ID not found in context"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var banana domain.BananaUpdate
	if err := json.NewDecoder(r.Body).Decode(&banana); err != nil {
		logError("update", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if banana.Name == "" || banana.Color == "" {
		logError("update", fmt.Errorf("incorrect request data"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.bananaService.Update(id, banana)
	if err != nil {
		logError("update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		logError("delete", fmt.Errorf("ID not found in context"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.bananaService.Delete(id)
	if err != nil {
		logError("delete", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
