CREATE TABLE teams (
	id SERIAL,
	name varchar(50) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE team_wishlist (
	teamID INT NOT NULL,
	profileID INT NOT NULL,
	PRIMARY KEY(teamID, profileID),
	CONSTRAINT fk_team
		FOREIGN KEY(teamID)
			REFERENCES teams(id)
);

CREATE TABLE participants (
	id SERIAL,
	name varchar(50) NOT NULL,
	teamID INT,
	PRIMARY KEY(id),
	CONSTRAINT fk_team
		FOREIGN KEY(teamID)
			REFERENCES teams(id)
);

CREATE TABLE profiles_participants (
	profileID INT NOT NULL,
	participantID INT NOT NULL,
	PRIMARY KEY(profileID, participantID),
	CONSTRAINT fk_participant
		FOREIGN KEY(participantID)
			REFERENCES participants(id)

);

CREATE TABLE contact_infos (
	email varchar(50) NOT NULL,
	participantID INT NOT NULL,
	CONSTRAINT fk_participant
		FOREIGN KEY(participantID)
			REFERENCES participants(id)
);
