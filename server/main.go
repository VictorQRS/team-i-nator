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

	router.HandleFunc("/participant", listParticipants).Methods("GET")
	router.HandleFunc("/participant/{id}", getParticipant).Methods("GET")
	router.HandleFunc("/participant/", newParticipant).Methods("POST")

	router.HandleFunc("/team/{id}", getTeam).Methods("GET")
	router.HandleFunc("/team/", newTeam).Methods("POST")
	router.HandleFunc("/team/{id}/participant", addParticipant).Methods("POST")
	router.HandleFunc("/team/{id}/wishlist", updateWishList).Methods("PUT")

	fmt.Println("Server at 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func listParticipants(writer http.ResponseWriter, req *http.Request) {
	participants, err := ListParticipants()
	if err != nil {
		// TODO: return error 500 message
	} else {
		json.NewEncoder(writer).Encode(participants)
	}
}

func getParticipant(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err);
		// TODO: show not found
	}

	participant, err := GetParticipant(id)
	if err != nil {
		// TODO: return error 500 message
	} else {
		json.NewEncoder(writer).Encode(participant)
	}
}

func newParticipant(writer http.ResponseWriter, req *http.Request) {
	var participant Participant
	_ = json.NewDecoder(req.Body).Decode(&participant)
	
	added, err := NewParticipant(participant)
	if err != nil {
		json.NewEncoder(writer).Encode(added)
	} else {
		// TODO: return error 500 message
	}
}

func getTeam(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		teams := ListTeams()
		json.NewEncoder(writer).Encode(teams)
		return
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err);
		// TODO: show not found
	}

	team := GetTeam(id)
	json.NewEncoder(writer).Encode(team);
}

func newTeam(writer http.ResponseWriter, req *http.Request) {
	var body struct{Name string `json:"name"`; FounderId int `json:"founderId"`}
	_ = json.NewDecoder(req.Body).Decode(&body)
	json.NewEncoder(writer).Encode(NewTeam(body.Name, body.FounderId))
}

func addParticipant(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err);
		// TODO: show not found
	}

	var body struct{ParticipantId int `json:"participantId"`}
	_ = json.NewDecoder(req.Body).Decode(&body)
	AddParticipant(id, body.ParticipantId)
}

func updateWishList(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err);
		// TODO: show not found
	}

	var body struct{Wishlist []Profile `json:"wishlist"`}
	_ = json.NewDecoder(req.Body).Decode(&body)
	UpdateWishList(id, body.Wishlist)
}
