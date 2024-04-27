package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"spacejunk3000/enemy"
	"spacejunk3000/player"
	"spacejunk3000/weapon"
	"strings"
	"time"
)

type Game struct {
	Player  *player.Player
	Enemies []enemy.Enemy
	Weapons []weapon.Weapon
}

// NewGame creates a new game instance or loads an existing one.
func NewGame(playerName string, charType player.CharacterType) (*Game, error) {
	// Load player data if it exists or create a new player.
	p, err := player.LoadPlayer(playerName)
	if err != nil || p == nil { // If there's no existing player, create a new one.
		p = &player.Player{
			Name:   playerName,
			Type:   charType,
			Health: 12, // Default health or other initialization parameters
			// Initialize other necessary fields
		}
		if err := player.SavePlayer(p); err != nil {
			return nil, fmt.Errorf("failed to save new player: %v", err)
		}
	}

	enemies, err := enemy.LoadEnemies("data/enemies.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load enemies: %v", err)
	}

	// If the player data was loaded successfully and the character type needs to be updated.
	if p != nil && p.Type != charType {
		p.Type = charType
		if err := player.SavePlayer(p); err != nil {
			return nil, fmt.Errorf("failed to update player type: %v", err)
		}
	}

	// Load weapons
	weapons, err := weapon.LoadWeapons("data/weapons.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load weapons: %v", err)
	}

	// Randomly select a weapon for the player
	randomIndex := rand.Intn(len(weapons))
	selectedWeapon := weapons[randomIndex]

	// Assign the selected weapon to the player
	p.Weapon = &selectedWeapon

	// Save the updated player data
	if err := player.SavePlayer(p); err != nil {
		return nil, fmt.Errorf("failed to save player: %v", err)
	}

	// Save weapons
	if err := weapon.SaveWeapons("data/weapons.json", weapons); err != nil {
		return nil, fmt.Errorf("failed to save weapons: %v", err)
	}

	return &Game{Player: p, Enemies: enemies, Weapons: weapons}, nil
}

// StartGame initializes and starts the game.
func StartGame(playerID string) (*player.Player, error) {
	p, err := player.LoadPlayer(playerID)
	if err != nil || p == nil { // New player or failed to load existing player
		charType := SelectCharacterType() // Let the user select a character type
		p = &player.Player{
			Name:   playerID,
			Type:   charType,
			Health: 12,
			Alive:  true,
		}
		if err := player.SavePlayer(p); err != nil {
			return nil, fmt.Errorf("failed to save new player: %v", err)
		}
	} else if !p.Alive { // Check if the player is starting over due to death
		player.ResetPlayer(p)
		player.SavePlayer(p)
	}

	return p, nil
}

// Start begins the game.
func (g *Game) Start() error {
	fmt.Println("Game has started.")
	fmt.Printf("Player %s has entered the game as a %s with %d health points and equipped with %s.\n", g.Player.Name, g.Player.Type, g.Player.Health, g.Player.Weapon.Name)

	if len(g.Enemies) == 0 {
		return fmt.Errorf("no enemies available for an encounter")
	}

	// Create a new local random generator
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Use the local random generator
	randomIndex := r.Intn(len(g.Enemies))
	randomEnemy := g.Enemies[randomIndex]

	fmt.Printf("A wild %s appears! It has %d health and can deal %d damage.\n", randomEnemy.Name, randomEnemy.Health, randomEnemy.Damage)

	// Implement encounter logic
	// if err := g.handleEncounter(randomEnemy); err != nil {
	//     return fmt.Errorf("encounter failed: %v", err)
	// }

	return nil
}

// UpdateHealth updates the player's health and handles death.
func UpdateHealth(p *player.Player, delta int) {
	p.Health += delta
	if p.Health <= 0 {
		p.Alive = false
		player.SavePlayer(p) // Save the dead state
		fmt.Println("You have died and must start over.")
		player.ResetPlayer(p) // Reset for new game start
	}
	player.SavePlayer(p) // Save any changes to player data
}

// SelectCharacterType prompts the user to select a character type.
func SelectCharacterType() player.CharacterType {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose your character type:")
	fmt.Println("1. Pirate")
	fmt.Println("2. Space Marine")
	fmt.Println("3. Empath")

	for {
		fmt.Print("Enter choice (1-3): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input, please try again.")
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1", "Pirate":
			return player.Pirate
		case "2", "Space Marine":
			return player.SpaceMarine
		case "3", "Empath":
			return player.Empath
		default:
			fmt.Println("Invalid choice, please select a valid character type.")
		}
	}
}
