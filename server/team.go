package main

import (
	"database/sql"
	"errors"
)

type Team struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type TeamFull struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Participants []ParticipantSimple `json:"participants"`
	WishList []Profile `json:"wishlist"`
}

func (t *TeamFull) Copy() TeamFull {
	return TeamFull {
		t.Id,
		t.Name,
		t.Participants,
		t.WishList,
	}
}

type TeamFullRow struct {
	Id int
	Name string
	ParticipantId int
	ParticipantName string
	ProfileId int
}

func GetTeam(id int) (TeamFull, error) {
	db := SetupDB()

	rows, err := db.Query(`
		SELECT
			teams.id, teams.name,
			participants.id, participants.name,
			team_wishlist.profileID
		FROM teams 
		LEFT OUTER JOIN participants ON teams.id = participants.teamID
		LEFT OUTER JOIN team_wishlist ON teams.id =  team_wishlist.teamID
		WHERE teams.id = $1;
	`, id)
	
	if err != nil {
		return TeamFull{}, err
	}
	
	var teamRows []TeamFullRow
	for rows.Next() {
		var (
			teamId int
			teamName string
			participantId int
			participantName string
			profileId sql.NullInt64
		)
		err = rows.Scan(&teamId, &teamName, &participantId, &participantName, &profileId)
		if err != nil {
			return TeamFull{}, err
		}

		_profileId := 0
		if (profileId.Valid) { _profileId = int(profileId.Int64) }

		teamRows = append(teamRows, TeamFullRow{
			teamId,
			teamName,
			participantId,
			participantName,
			_profileId,
		})
	}
	if len(teamRows) == 0 {
		return TeamFull{}, errors.New("no team")
	}

	teamsMap := make(map[int]*TeamFull)
	for _, row := range teamRows {
		if val, ok := teamsMap[row.Id]; ok {
			val.WishList = append(val.WishList, Profile(row.ProfileId))
			val.Participants = append(val.Participants, ParticipantSimple{
				row.ParticipantId,
				row.ParticipantName,
			})
		} else {
			teamsMap[row.Id] = &TeamFull{
				row.Id,
				row.Name,
				[]ParticipantSimple{{
					row.ParticipantId,
					row.ParticipantName,
				}},
				[]Profile{ Profile(row.ProfileId) },
			}
		}
	}

	teams := make([]TeamFull, 0, len(teamsMap))
	for  _, team := range teamsMap {
		teams = append(teams, team.Copy())
	}

	return teams[0], nil
}

func ListTeams() ([]Team, error) {
	db := SetupDB()

	rows, err := db.Query("SELECT teams.id, teams.name FROM teams;")
	if err != nil {
		return []Team{}, err
	}
	
	var teams []Team
	for rows.Next() {
		var (
			id int
			name string
		)

		err = rows.Scan(&id, &name)
		if err != nil {
			return []Team{}, err
		}

		teams = append(teams, Team{ id, name })
	}

	return teams, nil
}

func NewTeam(name string, founderId int) (TeamFull, error) {
	db := SetupDB()

	txn, err := db.Begin()
	if err != nil {
		return TeamFull{}, err
	}

	teamId, err := insertTeam(db, name);
	if err != nil {
		txn.Rollback()
		return TeamFull{}, err;
	}
	
	err = AssignToTeam(founderId, teamId, db)
	if err != nil {
		txn.Rollback()
		return TeamFull{}, err;
	}

	txn.Commit()

	return GetTeam(teamId)
}

func insertTeam(db *sql.DB, name string) (int,error) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO teams(name) VALUES($1) returning id;", name).Scan(&lastInsertID)
	return lastInsertID, err
}

func AddParticipant(teamId int, participantId int) error {
	db := SetupDB()

	err := AssignToTeam(participantId, teamId, db)
	if err != nil {
		return err;
	}

	return nil
}

func UpdateWishList(teamId int, wishlist []Profile) error {
	db := SetupDB()

	txn, err := db.Begin()
	if err != nil {
		return err
	}

	for _, profile := range wishlist {
		_, err = db.Exec("INSERT INTO team_wishlist(teamID, profileID) VALUES ($1, $2) ON CONFLICT DO UPDATE", teamId, profile);
		if err != nil {
			txn.Rollback()
			return err
		}
	}

	txn.Commit()
	return nil
}