package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikhailadvani/one-time-secret/pkg/api"
)

var getAPIBase = "/api/v1/secret"
var getAPI = fmt.Sprintf("%s/{secretID}", getAPIBase)
var createAPI = "/api/v1/secret"

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, api.Index())
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	secret, err := api.CreateSecret(r.Body, r.Host+getAPIBase)
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
	router.HandleFunc(createAPI, createSecret).Methods("POST")
	router.HandleFunc(getAPI, getSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
