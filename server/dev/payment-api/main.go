package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type PaymentRequest struct {
	Amount     int    `json:"amount"`
	ExternalID string `json:"external_id"`
}

type PaymentResponse struct {
	Data struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

var (
	store = make(map[string]string) // external_id => generated uuid
	mu    sync.Mutex
)

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.ExternalID == "" || req.Amount <= 0 {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// simulate slow response
	// time.Sleep(10 * time.Second)

	mu.Lock()
	defer mu.Unlock()

	id, exists := store[req.ExternalID]
	if !exists {
		id = uuid.NewString()
		store[req.ExternalID] = id

		resp := PaymentResponse{}
		resp.Data.ID = id
		resp.Data.ExternalID = req.ExternalID
		resp.Data.Status = "success"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// idempotency key conflict
	resp := PaymentResponse{}
	resp.Data.ID = id
	resp.Data.ExternalID = req.ExternalID
	resp.Data.Status = "success"
	resp.Message = "external id already exists"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := 9999
	http.HandleFunc("/v1/payments", paymentHandler)
	log.Printf("mock payment server running at port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
