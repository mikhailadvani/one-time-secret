package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
)

var getEndpointBase = "/api/v1/secret"
var getEndpoint = fmt.Sprintf("%s/{secretID}", getEndpointBase)
var createEndpoint = "/api/v1/secret"

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, api.Index(createEndpoint))
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var request api.CreateSecretRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "{}", http.StatusInternalServerError)
	}
	secret, err := api.CreateSecret(request, r.Host+getEndpointBase)
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
	router.HandleFunc(createEndpoint, createSecret).Methods("POST")
	router.HandleFunc(getEndpoint, getSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
