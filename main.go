package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Email struct {
	Assunto  string `json:"assunto"`
	Mensagem string `json:"mensage"`
}

var emails = map[string]*Email{
	"1": &Email{Assunto: "Contato newsletter", Mensagem: "Olá paulo como está"},
	"2": &Email{Assunto: "Contato newsletter", Mensagem: "Olá joão como está"},
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/emails", handleEmails)
	router.HandleFunc("/email/{id}", handleEmail).Methods("GET", "DELETE")

	http.ListenAndServe(":8080", router)
}

func handleEmail(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	id := vars["id"]

	switch request.Method {
	case "GET":
		email, ok := emails[id]
		if !ok {
			response.WriteHeader(http.StatusNotFound)
			fmt.Fprint(response, string("E-mail not found"))
		}
		outgoingJSON, error := json.Marshal(email)
		if error != nil {
			log.Println(error.Error())
			http.Error(response, error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(response, string(outgoingJSON))
	case "DELETE":
		delete(emails, id)
		response.WriteHeader(http.StatusAccepted)
	case "POST":
		email := new(Email)
		decoder := json.NewDecoder(request.Body)
		error := decoder.Decode(&email)
		if error != nil {
			log.Print(error.Error())
			http.Error(response, error.Error(), http.StatusInternalServerError)
			return
		}
		emails[id] = email
		outgoingJSON, error := json.Marshal(email)
		if error != nil {
			log.Print(error.Error())
			http.Error(response, error.Error(), http.StatusInternalServerError)
			return
		}
		response.WriteHeader(http.StatusCreated)
		fmt.Fprint(response, string(outgoingJSON))
	}
}

func handleEmails(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	outgoingJSON, error := json.Marshal(emails)
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}
