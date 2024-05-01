package main

import (
	"flag"
	"fmt"
	"log"
	"spacejunk3000/doorutil"
	"spacejunk3000/enemy"
	"spacejunk3000/game"
	"spacejunk3000/implant"
	"spacejunk3000/location"
	"spacejunk3000/player"
	"spacejunk3000/weapon"
)

func main() {
	doorutil.ClearScreen()

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

		doorutil.ClearScreen()
		doorutil.CursorHide()
		doorutil.DisplayAnsiFile("assets/start.ans", true)

		err := doorutil.WaitForAnyKey()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Select character type
		charType := game.SelectCharacterType()

		// Load implants from JSON file
		implants, err := implant.LoadImplants("data/implants.json")
		if err != nil {
			log.Fatalf("Failed to load implants: %v", err)
		}

		// Select implant
		selectedImplant := implant.SelectImplant(implants)

		// Create a new player with default values, dropfile information, character type, and selected implant
		p, err = player.NewPlayer(playerName, charType, dropTimeLeft, nodeNum, dropEmulation)
		if err != nil {
			log.Fatalf("Failed to create new player: %v", err)
		}

		// Set the selected implant for the player
		p.Implant = selectedImplant

		// Save the new player
		if err := player.SavePlayer(p); err != nil {
			log.Fatalf("Failed to save new player: %v", err)
		}
	}

	// Load locations from JSON file
	locations, err := location.LoadLocations("data/locations.json")
	if err != nil {
		log.Fatalf("Failed to load locations: %v", err)
	}

	// Load enemies from JSON file
	enemies, err := enemy.LoadEnemies("data/enemies.json")
	if err != nil {
		log.Fatalf("Failed to load enemies: %v", err)
	}

	// Load weapons from JSON file
	weapons, err := weapon.LoadWeapons("data/weapons.json")
	if err != nil {
		log.Fatalf("Failed to load weapons: %v", err)
	}

	// Load implants from JSON file
	implants, err := implant.LoadImplants("data/implants.json")
	if err != nil {
		log.Fatalf("Failed to load implants: %v", err)
	}

	// Initialize and start the game with all required arguments
	g, err := game.NewGame(playerName, p.Type, weapons, implants, locations, enemies)
	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	// Start the encounter
	game.StartNewEncounter(g)
}
