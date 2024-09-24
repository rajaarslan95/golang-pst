package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"user-manager/schemas"
	"user-manager/services"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc services.UserService
	ch  chan string
	mu  *sync.RWMutex
	wg  *sync.WaitGroup
}

func NewUserHandler(svc services.UserService, ch chan string, wg *sync.WaitGroup) *Handler {
	return &Handler{
		svc: svc,
		ch:  ch,
		mu:  &sync.RWMutex{},
		wg:  wg,
	}
}

// Function to log request duration
func (h *Handler) logDuration(id string, start time.Time) {
	defer h.wg.Done()
	duration := time.Since(start)
	h.ch <- fmt.Sprintf("Request %s took %v", id, duration)
}

func (h *Handler) GetSvc() services.UserService {
	return h.svc
}

// Create a user
func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	start := time.Now()
	defer h.logDuration("addUser", start)

	h.mu.Lock() // Ensure consistent writes
	defer h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	var user schemas.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Unable to parse body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schemas.NewError("Invalid input", http.StatusBadRequest))
		return
	}

	if err := h.svc.Create(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(schemas.NewError("Server failure", http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schemas.NewError("User Created", http.StatusCreated))

}

// Update a user
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	start := time.Now()
	defer h.logDuration("UpdateUser", start)

	h.mu.Lock() // Ensure consistent writes
	defer h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schemas.NewError("Invalid user ID", http.StatusBadRequest))
		return
	}

	var user schemas.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Unable to parse body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schemas.NewError("Invalid input", http.StatusBadRequest))
		return
	}
	user.ID = id
	if err := h.svc.Update(user); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schemas.NewError("Not found", http.StatusNotFound))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schemas.NewError("User Updated", http.StatusOK))
}

// Delete a user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	start := time.Now()
	defer h.logDuration("DeleteUser", start)

	h.mu.Lock() // Ensure consistent writes
	defer h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schemas.NewError("Invalid user ID", http.StatusBadRequest))
		return
	}

	if err := h.svc.Delete(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schemas.NewError("Not found", http.StatusNotFound))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schemas.NewError("User Deleted", http.StatusOK))
}

// Get a user
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	start := time.Now()
	defer h.logDuration("GetUser", start)

	h.mu.Lock() // Ensure consistent writes
	defer h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schemas.NewError("Invalid user ID", http.StatusBadRequest))
		return
	}

	user, err := h.svc.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schemas.NewError("Not found", http.StatusNotFound))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
