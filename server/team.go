package main

type Team struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type TeamFull struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Participants []Participant `json:"participants"`
	WishList []Profile `json:"wishlist"`
}

func GetTeam(id int) TeamFull {
	return TeamFull{}
}

func ListTeams() []Team {
	return []Team{}
}

func NewTeam(name string, founderId int) TeamFull {
	participants := make([]Participant, 1)
	participants[0] = GetParticipant(founderId)

	return TeamFull {
		0,
		name,
		participants,
		make([]Profile, 0),
	}
}

func AddParticipant(teamId int, participantId int) {
	// set teamId of Partificipant of id `id`
}

func UpdateWishList(teamId int, wishlist []Profile) {
	// upsert tuple (teamId, profileId)
}