package player

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taml5/team-matchmaker.git/internal/config"
)

var Positions = []string{"top", "jgl", "mid", "adc", "sup"}

type Player struct {
	UUID     string // UUID is the unique identifier for the player
	Username string
	Tag      string
	Name     string
	Elo      map[string]int
}

type Team map[string]Player

type RiotAccount struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

func NewPlayer(name string, username string, tag string) (*Player, error) {
	uuid, err := getID(username, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get ID for player %s#%s: %v",
			username,
			tag,
			err)
	}

	// Default elo is 1000 for all roles
	defaultElo := map[string]int{
		"top": 1000,
		"jgl": 1000,
		"mid": 1000,
		"adc": 1000,
		"sup": 1000,
	}

	newPlayer := &Player{
		UUID:     uuid,
		Username: username,
		Tag:      tag,
		Name:     name,
		Elo:      defaultElo,
	}

	err = InsertPlayer(newPlayer)
	if err != nil {
		return nil, fmt.Errorf("failed to insert player into database: %v", err)
	}
	return newPlayer, nil
}

func getID(username string, tag string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf(
		"https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		username,
		tag,
	)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", config.RiotAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var account RiotAccount
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return account.Puuid, nil
}
