package player

import (
	"fmt"

	"github.com/taml5/team-matchmaker.git/internal/config"
)

func GetPlayerByID(uuid string) (*Player, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("GetPlayerByID failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT uuid, username, tag, name FROM player WHERE uuid = ?`,
		uuid,
	)

	var player Player
	err = row.Scan(&player.UUID, &player.Username, &player.Tag, &player.Name)
	if err != nil {
		return nil, fmt.Errorf("GetPlayerByID failed: %w", err)
	}

	player.Elo = make(map[string]int)
	for _, position := range Positions {
		var elo int
		err = tx.QueryRow(
			`SELECT elo FROM elo WHERE player_id = ? AND position = ?`,
			player.UUID,
			position,
		).Scan(&elo)
		if err != nil {
			return nil, fmt.Errorf("failed to get elo for position %s: %w", position, err)
		}
		player.Elo[position] = elo
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &player, nil
}

func GetPlayerByUsernameAndTag(username string, tag string) (*Player, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("GetPlayerByUsernameAndTag failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT uuid, username, tag, name FROM player WHERE username = ? AND tag = ?`,
		username,
		tag,
	)

	var player Player
	err = row.Scan(&player.UUID, &player.Username, &player.Tag, &player.Name)
	if err != nil {
		return nil, fmt.Errorf("GetPlayerByUsernameAndTag failed: %w", err)
	}

	player.Elo = make(map[string]int)
	for _, position := range Positions {
		var elo int
		err = tx.QueryRow(
			`SELECT elo FROM elo WHERE player_id = ? AND position = ?`,
			player.UUID,
			position,
		).Scan(&elo)
		if err != nil {
			return nil, fmt.Errorf("failed to get elo for position %s: %w", position, err)
		}
		player.Elo[position] = elo
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &player, nil
}

func InsertPlayer(player *Player) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("insertPlayer failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO Player (uuid, username, tag, name) VALUES (?, ?, ?, ?)`,
		player.UUID,
		player.Username,
		player.Tag,
		player.Name,
	)
	if err != nil {
		return fmt.Errorf("insertPlayer failed: %w", err)
	}

	// insert into elo tables
	for _, position := range Positions {
		_, err = tx.Exec(
			`INSERT INTO Elo (player_id, position, elo) VALUES (?, ?, ?)`,
			player.UUID,
			position,
			player.Elo[position],
		)
		if err != nil {
			return fmt.Errorf("insertPlayer failed for position %s: %w",
				position,
				err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func UpdateRoleElo(player Player, role string, newElo int) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("UpdateRoleElo failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`UPDATE Elo SET elo = ? WHERE player_id = ? AND position = ?`,
		newElo,
		player.UUID,
		role,
	)
	if err != nil {
		return fmt.Errorf("UpdateRoleElo failed for role %s: %w", role, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
