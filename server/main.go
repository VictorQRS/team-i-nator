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

	router.HandleFunc("/participant/", listParticipants).Methods("GET")
	router.HandleFunc("/participant/{id}", getParticipant).Methods("GET")
	router.HandleFunc("/participant/", newParticipant).Methods("POST")

	router.HandleFunc("/team/", listTeams).Methods("GET")
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
		fmt.Println(err);
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
		fmt.Println(err);
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
		fmt.Println(err);
		// TODO: return error 500 message
	} else {
		json.NewEncoder(writer).Encode(added)
	}
}

func listTeams(writer http.ResponseWriter, req *http.Request) {
	teams, err := ListTeams()
	if err != nil {
		fmt.Println(err);
		// TODO: return error 500 message
	} else {
		json.NewEncoder(writer).Encode(teams)
	}
}

func getTeam(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err);
		// TODO: show not found
	}

	team, err := GetTeam(id)
	if err != nil {
		fmt.Println(err);
		// TODO: return error 500 message
	} else {
		json.NewEncoder(writer).Encode(team)
	}
}

func newTeam(writer http.ResponseWriter, req *http.Request) {
	var body struct{Name string `json:"name"`; FounderId int `json:"founderId"`}
	_ = json.NewDecoder(req.Body).Decode(&body)
	team, err := NewTeam(body.Name, body.FounderId)
	if err != nil {
		fmt.Println(err);
		// TODO: return error 500 message
	} else {
		json.NewEncoder(writer).Encode(team)
	}
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
	err = AddParticipant(id, body.ParticipantId)
	if err != nil {
		fmt.Println(err);
		// TODO: show 500
	}
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
	err = UpdateWishList(id, body.Wishlist)
	if err != nil {
		fmt.Println(err);
		// TODO: show 500
	}
}
