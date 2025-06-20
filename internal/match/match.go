package match

import (
	"math"

	"github.com/taml5/team-matchmaker.git/internal/player"
)

type Match struct {
	Team1  player.Team
	Team2  player.Team
	Result bool // true for Team1 win, false for Team2 win
}

// This function will update the Elos of both teams based on the match result.
func UpdateElos(match *Match) error {
	var player1, player2 player.Player
	var newElo1, newElo2 int
	teamEloChange1, teamEloChange2 := teamEloChange(
		match.Team1,
		match.Team2,
		match.Result,
	)

	for _, pos := range player.Positions {
		player1 = match.Team1[pos]
		player2 = match.Team2[pos]

		personalEloChange1, personalEloChange2 := personalEloChange(
			player1,
			player2,
			pos,
			match.Result,
		)

		// Slightly weighted towards personal Elo change
		eloChange1 := 0.6*float64(personalEloChange1) + 0.4*float64(teamEloChange1)
		eloChange2 := 0.6*float64(personalEloChange2) + 0.4*float64(teamEloChange2)
		newElo1 = player1.Elo[pos] + int(eloChange1)
		newElo2 = player2.Elo[pos] + int(eloChange2)

		// Update the Elos of both players
		player.UpdateRoleElo(player1, pos, newElo1)
		player.UpdateRoleElo(player2, pos, newElo2)
	}
	return nil
}

func meanElo(team player.Team) int {
	totalElo := 0
	for _, pos := range player.Positions {
		totalElo += team[pos].Elo[pos]
	}
	return totalElo / len(team)
}

// Return the personal ELO change for player1 and player2, based on the
// difference in their ELOs.
func personalEloChange(
	player1 player.Player,
	player2 player.Player,
	position string,
	result bool, // true if player1 wins, false if player2 wins
) (int, int) {
	elo1, elo2 := player1.Elo[position], player2.Elo[position]
	return getEloChange(elo1, elo2, result)
}

// Return the mean team ELO change between each team, based on the
// mean ELO difference in their ELOs and the match result.
func teamEloChange(
	team1 player.Team,
	team2 player.Team,
	result bool, // true if team1 wins, false if team2 wins
) (int, int) {
	elo1, elo2 := meanElo(team1), meanElo(team2)
	return getEloChange(elo1, elo2, result)
}

func getEloChange(
	elo1 int,
	elo2 int,
	result bool,
) (int, int) {
	k := 32.0
	expectedScore1 := 1 / (1 + math.Pow(10, float64(elo2-elo1)/400.0))
	expectedScore2 := 1 / (1 + math.Pow(10, float64(elo1-elo2)/400.0))

	var change1, change2 float64
	if result {
		change1 = k * (1.0 - expectedScore1)
		change2 = k * (0.0 - expectedScore2)
	} else {
		change1 = k * (0.0 - expectedScore1)
		change2 = k * (1.0 - expectedScore2)
	}
	return int(change1), int(change2)
}
