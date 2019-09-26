package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/mux"
	"github.com/edcast/one-time-secret/pkg/api"
	configf "github.com/edcast/one-time-secret/pkg/config"
)

var config = configf.LoadConfig()
var getEndpointBase = "/api/secret"
var getEndpoint = fmt.Sprintf("%s/{secretID}", getEndpointBase)
var createEndpoint = "/api/secret"

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
	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		http.Error(w, "{}", http.StatusInternalServerError)
	}
	baseURL.Path = path.Join(baseURL.Path, getEndpointBase)
	urlString := baseURL.String()
	secret, err := api.CreateSecret(request, urlString)
	if err != nil {
		http.Error(w, "{}", secret.Status)
	}
	json.NewEncoder(w).Encode(secret)
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
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
