package main

type Participant struct {
	Name string
	Profile []Profile
	ContactInfo ContactInfo
}

func NewParticipant(name string, profile []Profile, info ContactInfo) Participant {
	return Participant {
		name,
		profile,
		info,
	}
}