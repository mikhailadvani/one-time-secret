package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, api.Index())
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	secret, err := api.CreateSecret(r.Body)
	if err != nil {
		http.Error(w, "{}", secret.Status)
	}
	json.NewEncoder(w).Encode(secret)
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secretID := vars["secretID"]
	secret, err := api.GetSecret(secretID)
	if err != nil {
		http.Error(w, "{}", secret.Status)
	}
	fmt.Fprintln(w, secret.Content)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/api/v1/secret", createSecret).Methods("POST")
	router.HandleFunc("/api/v1/secret/{secretID}", getSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
