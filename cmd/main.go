package main

import (
	"fmt"

	"github.com/taml5/team-matchmaker.git/internal/config"
	"github.com/taml5/team-matchmaker.git/internal/player"
)

func main() {
	config.LoadEnv()
	err := config.LoadDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to load database: %v", err))
	}

	newPlayer, err := player.NewPlayer("Dexter", "fhweo", "fhweo")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Player created: %+v\n", newPlayer)
}
