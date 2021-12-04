package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/participant/{id}", getParticipant).Methods("GET")
	router.HandleFunc("/participant/", newParticipant).Methods("POST")

	fmt.Println("Server at 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func getParticipant(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		participants := ListParticipants()
		json.NewEncoder(writer).Encode(participants)
		return
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err);
		// show not found
	}

	participant := GetParticipant(id)
	json.NewEncoder(writer).Encode(participant);
}

func newParticipant(writer http.ResponseWriter, req *http.Request) {
	var participant Participant
	_ = json.NewDecoder(req.Body).Decode(&participant)
	json.NewEncoder(writer).Encode(NewParticipant(participant))
}
