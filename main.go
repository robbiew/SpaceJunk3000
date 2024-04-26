package main

import (
	"flag"
	"fmt"
	"log"
	"spacejunk3000/game"
	"spacejunk3000/player"
	"spacejunk3000/util"
)

func main() {
	// Define flags
	dropfilePath := flag.String("dropfile", "", "path to the BBS drop file")
	flag.Parse()

	// Check if dropfile flag is provided
	if *dropfilePath == "" {
		log.Fatal("Dropfile path is required. Please provide the path using the -dropfile flag.")
	}

	// Get BBS dropfile information about the user using the DropFileData function
	dropAlias, dropTimeLeft, dropEmulation, nodeNum := util.DropFileData(*dropfilePath)

	// Use dropAlias as the playerName
	playerName := dropAlias

	// Load or create player
	p, err := player.LoadPlayer(playerName)
	if err != nil {
		fmt.Println("No existing player found or error loading player:", err)
		fmt.Println("Please select your character type:")
		charType := game.SelectCharacterType()

		// Create a new player with default values and dropfile information
		p = &player.Player{
			Name:      playerName,
			Type:      charType,
			Health:    dropTimeLeft,                                   // Using dropTimeLeft as health
			Inventory: []player.ItemType{player.Sword, player.Shield}, // Initialize inventory here
			Alive:     true,                                           // Set other necessary fields
			Stats: player.Stats{
				Might:   dropEmulation, // Using dropEmulation as Might
				Cunning: 2,
				Wisdom:  nodeNum, // Using nodeNum as Wisdom
			},
		}

		// Save the new player
		if err := player.SavePlayer(p); err != nil {
			log.Fatalf("Failed to save new player: %v", err)
		}
	}

	// Initialize and start the game
	g, err := game.NewGame(playerName, p.Type) // Pass playerName and player type
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	err = g.Start()
	if err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}
}
