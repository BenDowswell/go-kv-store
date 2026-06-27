package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BenDowswell/go-kv-store/kvstore"
)

type SetRequest struct {
	Value string `json:"value"`
}

func Start(store *kvstore.KVStore) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", getHealth)
	mux.HandleFunc("GET /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		getKey(w, r, store)
	})
	mux.HandleFunc("PUT /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		updateKey(w, r, store)
	})
	mux.HandleFunc("DELETE /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		deleteKey(w, r, store)
	})

	fmt.Println("HTTP Server Listening on 8080")
	http.ListenAndServe(":8080", mux)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func getKey(w http.ResponseWriter, r *http.Request, store *kvstore.KVStore) {
	key := r.PathValue("key")

	value, ok := store.Get(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, value)
}

func updateKey(w http.ResponseWriter, r *http.Request, store *kvstore.KVStore) {
	key := r.PathValue("key")

	var req SetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Value == "" {
		http.Error(w, "value is required", http.StatusBadRequest)
		return
	}

	store.Set(key, req.Value)
	w.WriteHeader(http.StatusNoContent)
}

func deleteKey(w http.ResponseWriter, r *http.Request, store *kvstore.KVStore) {
	key := r.PathValue("key")

	store.Delete(key)
	w.WriteHeader(http.StatusNoContent)
}
