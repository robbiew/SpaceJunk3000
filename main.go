package main

import (
	"flag"
	"fmt"
	"log"
	"spacejunk3000/doorutil"
	"spacejunk3000/game"
	"spacejunk3000/player"
	"spacejunk3000/weapon"
)

func main() {
	// Define flags
	dropfilePath := flag.String("door32", "", "path to the Door32.sys drop file")
	flag.Parse()

	// Check if dropfile flag is provided
	if *dropfilePath == "" {
		log.Fatal("Dropfile path is required. Please provide the path using the -door32 flag.")
	}

	// Get BBS dropfile information about the user
	dropAlias, dropTimeLeft, dropEmulation, nodeNum := doorutil.DropFileData(*dropfilePath)

	// Use dropAlias as the playerName
	playerName := dropAlias

	// Load or create player
	p, err := player.LoadPlayer(playerName)
	if err != nil {
		fmt.Println("Welcome, New Player! Let's get you set up.")
		fmt.Println("Please select your character type:")
		charType := game.SelectCharacterType()

		// Create a new player with default values and dropfile information
		p, err = player.NewPlayer(playerName, charType, dropTimeLeft, nodeNum, dropEmulation)
		if err != nil {
			log.Fatalf("Failed to create new player: %v", err)
		}

		// Save the new player
		if err := player.SavePlayer(p); err != nil {
			log.Fatalf("Failed to save new player: %v", err)
		}
	}

	// Load weapons
	weapons, err := weapon.LoadWeapons("data/weapons.json")
	if err != nil {
		log.Fatalf("Failed to load weapons: %v", err)
	}

	// Initialize and start the game
	g, err := game.NewGame(playerName, p.Type, weapons)
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// Pass the list of weapons to the StartGame function
	player, err := game.StartGame(playerName, g.Weapons)
	if err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}

	fmt.Printf("Player %s has entered the game as a %s with %d health points", player.Name, player.Type, player.Health)

	// Check if the player has any weapons equipped
	if len(player.Weapons) > 0 {
		fmt.Printf(" and is equipped with:\n")
		for _, w := range player.Weapons {
			fmt.Printf("- %s\n", w.Name)
		}
	} else {
		fmt.Println(" and is unarmed.")
	}
}
