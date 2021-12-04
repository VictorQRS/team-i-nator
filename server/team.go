package main

type Team struct {
	Name string
	Participants []Participant
	WishList []Profile
}

func NewTeam(name string, participant Participant) Team {
	participants := make([]Participant, 1)
	participants[0] = participant
	
	return Team {
		name,
		participants,
		make([]Profile, 0),
	}
}

func (team *Team) AddParticipant(participant Participant) {
	team.Participants = append(team.Participants, participant)
}

func (team *Team) UpdateWishList(wishlist []Profile) {
	team.WishList = wishlist
}