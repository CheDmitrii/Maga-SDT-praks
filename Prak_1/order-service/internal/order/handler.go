package order

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repo   *Repo
	client *UserServiceClient
}

func NewHandler(repo *Repo, client *UserServiceClient) *Handler {
	return &Handler{
		repo:   repo,
		client: client,
	}
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "order id is required", http.StatusBadRequest)
		return
	}

	// /orders/{id}/full should not enter here
	if strings.HasSuffix(path, "/full") {
		http.Error(w, "use /orders/{id}/full handler", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrderWithUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected path: /orders/{id}/full
	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	if !strings.HasSuffix(path, "/full") {
		http.Error(w, "bad route", http.StatusBadRequest)
		return
	}

	idPart := strings.TrimSuffix(path, "/full")
	idPart = strings.TrimSuffix(idPart, "/")

	id, err := strconv.ParseInt(idPart, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	user, err := h.client.GetUserByID(order.UserID)
	if err != nil {
		http.Error(w, "failed to get user data: "+err.Error(), http.StatusBadGateway)
		return
	}

	result := OrderWithUser{
		Order: order,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}
