package main

import (
	"fmt"
	"log"
	"spacejunk3000/game"
	"spacejunk3000/player"
)

func main() {
	fmt.Println("Welcome to Space Junk 3000")

	// Assuming playerName is retrieved dynamically from BBS drop file
	playerName := "ExamplePlayerName" // Placeholder: replace with actual code to retrieve from BBS drop file

	// Load or create player
	p, err := player.LoadPlayer(playerName)
	if err != nil {
		fmt.Println("No existing player found or error loading player:", err)
		fmt.Println("Please select your character type:")
		charType := game.SelectCharacterType()

		// Create a new player with default values
		p = &player.Player{
			Name:      playerName,
			Type:      charType,
			Health:    12,
			Inventory: []player.ItemType{player.Sword, player.Shield}, // Initialize inventory here
			Alive:     true,                                           // Set other necessary fields
			Stats: player.Stats{
				Might:   1,
				Cunning: 2,
				Wisdom:  4,
			},
		}

		// Save the new player
		if err := player.SavePlayer(p); err != nil {
			log.Fatalf("Failed to save new player: %v", err)
		}
	}
}
