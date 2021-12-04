package main

type Participant struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Profile []Profile `json:"profiles"`
	ContactInfo ContactInfo `json:"contactInfo"`
}

func GetParticipant(id int) Participant {
	return Participant{}
}

func ListParticipants() []Participant {
	return []Participant{}
}

func NewParticipant(participant Participant) Participant {
	return participant
}