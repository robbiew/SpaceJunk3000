package main

import (
	"flag"
	"fmt"
	"log"
	"spacejunk3000/doorutil"
	"spacejunk3000/game"
	"spacejunk3000/player"
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
	dropAlias, dropTimeLeft, dropEmulation, nodeNum := doorutil.DropFileData(*dropfilePath)

	// Use dropAlias as the playerName
	playerName := dropAlias

	// Load or create player
	p, err := player.LoadPlayer(playerName)
	if err != nil {
		// fmt.Println("No existing player found or error loading player:", err)
		fmt.Println("Please select your character type:")
		charType := game.SelectCharacterType()

		// Create a new player with default values and dropfile information
		p = &player.Player{
			Name:   playerName,
			Type:   charType,
			Health: 12,
			Alive:  true,
			Stats: player.Stats{
				Might:   4,
				Cunning: 2,
				Wisdom:  1,
			},
			TimeLeft:  dropTimeLeft,
			Emulation: dropEmulation,
			NodeNum:   nodeNum,
			MaxSlots:  4,
		}

		// Save the new player
		if err := player.SavePlayer(p); err != nil {
			log.Fatalf("Failed to save new player: %v", err)
		}
	}

	// Initialize and start the game
	g, err := game.NewGame(playerName, p.Type)
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// Pass the list of weapons to the StartGame function
	player, err := game.StartGame(playerName, g.Weapons)
	if err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}

	fmt.Printf("Player %s has entered the game as a %s with %d health points", player.Name, player.Type, player.Health)

	// Check if the player has a weapon equipped
	if player.Weapon != nil {
		fmt.Printf(" and is equipped with %s.\n", player.Weapon.Name)
	} else {
		fmt.Println(" and is unarmed.")
	}
}
