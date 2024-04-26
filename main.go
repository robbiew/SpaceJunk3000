package main

import (
	"fmt"
	"log"
	"spacejunk3000/game"
	"spacejunk3000/player"
)

func main() {
	fmt.Println("Welcome to Space Junk 3000")
	playerName := "EnterYourName" // Placeholder for dynamic player name retrieval
	charType := player.Pirate     // Placeholder for dynamic character type selection

	g, err := game.NewGame(playerName, charType)
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	err = g.Start()
	if err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}
}
