package main

import (
	"database/sql"
)

type ParticipantSimple struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Participant struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Team Team `json:"team"`
	Profiles []Profile `json:"profiles"`
	ContactInfo ContactInfo `json:"contactInfo"`
}

func (p *Participant) Copy() Participant {
	return Participant{
		p.Id,
		p.Name,
		p.Team,
		p.Profiles,
		p.ContactInfo,
	}
}

type ParticipantRow struct {
	Id int
	Name string
	TeamId int
	TeamName string
	ProfileId int
	Email string
}

func getParticipants(id int) ([]Participant, error) {
	db := SetupDB()

	var rows *sql.Rows
	var err error
	if (id == 0) {
		rows, err = db.Query(`
			SELECT
				participants.id, participants.name,
				teams.id, teams.name,
				profiles_participants.profileID,
				contact_infos.email
			FROM participants 
			LEFT OUTER JOIN teams ON teams.id = participants.teamID
			LEFT OUTER JOIN profiles_participants ON profiles_participants.participantID = participants.id 
			LEFT OUTER JOIN contact_infos ON contact_infos.participantID = participants.id;
		`)
	} else {
		rows, err = db.Query(`
			SELECT
				participants.id, participants.name,
				teams.id, teams.name,
				profiles_participants.profileID,
				contact_infos.email
			FROM participants 
			LEFT OUTER JOIN teams ON teams.id = participants.teamID
			LEFT OUTER JOIN profiles_participants ON profiles_participants.participantID = participants.id 
			LEFT OUTER JOIN contact_infos ON contact_infos.participantID = participants.id
			WHERE participants.id = $1;
		`, id)
	}
	
	if err != nil {
		return []Participant{}, err
	}
	
	var participantRows []ParticipantRow
	for rows.Next() {
		var (
			id int
			name string
			teamId sql.NullInt64
			teamName sql.NullString
			profileId int
			email string
		)

		err = rows.Scan(&id, &name, &teamId, &teamName, &profileId, &email)
		if err != nil {
			return []Participant{}, err
		}

		_teamId := 0
		if (teamId.Valid) { _teamId = int(teamId.Int64) }
		_teamName := ""
		if (teamName.Valid) { _teamName = teamName.String }

		participantRows = append(participantRows, ParticipantRow{
			id,
			name,
			_teamId,
			_teamName,
			profileId,
			email,
		})
	}

	participantsMap := make(map[int]*Participant)
	for _, row := range participantRows {
		if val, ok := participantsMap[row.Id]; ok {
			val.Profiles = append(val.Profiles, Profile(row.ProfileId))
		} else {
			participantsMap[row.Id] = &Participant{
				row.Id,
				row.Name,
				Team { row.TeamId, row.TeamName },
				[]Profile{ Profile(row.ProfileId) },
				ContactInfo{ row.Email },
			}
		}
	}

	participants := make([]Participant, 0, len(participantsMap))
	for  _, participant := range participantsMap {
		participants = append(participants, participant.Copy())
	}

	return participants, nil
}

func GetParticipant(id int) (Participant, error) {
	participants, err := getParticipants(id);
	if err == nil || len(participants) > 0 {
		return participants[0], err
	}

	return Participant{}, err
}

func ListParticipants() ([]Participant, error) {
	return getParticipants(0);
}

func NewParticipant(participant Participant) (Participant, error) {
	db := SetupDB()

	txn, err := db.Begin()
	if err != nil {
		return Participant{}, err
	}

	participant.Id, err = insertParticipant(db, participant.Name);
	if err != nil {
		txn.Rollback()
		return Participant{}, err;
	}
	
	err = insertParticipantProfiles(db, participant.Id, participant.Profiles)
	if err != nil {
		txn.Rollback()
		return Participant{}, err;
	}
	
	err = insertContactInfo(db, participant.Id, participant.ContactInfo)
	if err != nil {
		txn.Rollback()
		return Participant{}, err;
	}

	txn.Commit()

	return participant, nil
}

func insertParticipant(db *sql.DB, name string) (int,error) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO participants(name) VALUES($1) returning id;", name).Scan(&lastInsertID)
	return lastInsertID, err
}

func insertParticipantProfiles(db *sql.DB, participantId int, profiles []Profile) error {
	for _, profile := range profiles {
		_, err := db.Exec("INSERT INTO profiles_participants(profileID, participantID) VALUES($1, $2);", profile, participantId)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertContactInfo(db *sql.DB, participantId int, contactInfo ContactInfo) error {
	_, err := db.Exec("INSERT INTO contact_infos(email, participantID) VALUES($1, $2);", contactInfo.Email, participantId)
	if err != nil {
		return err
	}
	return nil
}

func AssignToTeam(id int, teamId int, db *sql.DB) error {
	if (db == nil) {
		db = SetupDB()
	}

	// this does not stop if user does not exist, bad for team creation
	_, err := db.Exec("UPDATE participants SET teamID = $1 WHERE id = $2;", teamId, id)
	if err != nil {
		return err
	}
	return nil

}