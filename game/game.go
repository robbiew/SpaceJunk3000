package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"spacejunk3000/enemy"
	"spacejunk3000/item"
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

// InitializePlayer initializes a player by loading an existing one or creating a new one.
func InitializePlayer(playerName string, weapons []weapon.Weapon) (*player.Player, error) {
	// Load existing player or create a new one if not found
	p, err := player.LoadPlayer(playerName)
	if err != nil || p == nil {
		charType := SelectCharacterType() // Let the user select a character type if creating a new player
		p, err = player.NewPlayer(playerName, charType, 0, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to create new player: %v", err)
		}

		// Randomly select a weapon for the player
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		randomIndex := random.Intn(len(weapons))

		// Equip the randomly selected weapon to the player
		if err := EquipWeapon(p, &weapons[randomIndex]); err != nil {
			return nil, fmt.Errorf("failed to equip weapon: %v", err)
		}

		// Save the player data
		if err := player.SavePlayer(p); err != nil {
			return nil, fmt.Errorf("failed to save new player: %v", err)
		}

	} else if !p.Alive { // Check if the player is starting over due to death
		player.ResetPlayer(p)
		player.SavePlayer(p)
	} else if p.Weapons == nil { // Check if the player does not have a weapon equipped
		// Randomly select a weapon for the player
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		randomIndex := random.Intn(len(weapons))

		// Equip the randomly selected weapon to the player
		if err := EquipWeapon(p, &weapons[randomIndex]); err != nil {
			return nil, fmt.Errorf("failed to equip weapon: %v", err)
		}

		// Save the player data
		if err := player.SavePlayer(p); err != nil {
			return nil, fmt.Errorf("failed to save player: %v", err)
		}
	}

	return p, nil
}

// NewGame creates a new game instance or loads an existing one.
func NewGame(playerName string, charType player.CharacterType, weapons []weapon.Weapon) (*Game, error) {
	// Initialize the player
	p, err := InitializePlayer(playerName, weapons)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize player: %v", err)
	}

	// Load enemies
	enemies, err := enemy.LoadEnemies("data/enemies.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load enemies: %v", err)
	}

	return &Game{Player: p, Enemies: enemies}, nil
}

// StartGame initializes and starts the game.
func StartGame(playerName string, weapons []weapon.Weapon) (*player.Player, error) {
	// Initialize the player
	p, err := InitializePlayer(playerName, weapons)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize player: %v", err)
	}

	return p, nil
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

// EquipWeapon equips a weapon to the player if there are available slots.
func EquipWeapon(p *player.Player, w *weapon.Weapon) error {
	// Check if there are enough weapon slots to equip the weapon
	if p.WeaponSlots+w.Slots > p.MaxSlots {
		return fmt.Errorf("cannot carry that much")
	}

	// Equip the weapon to the player
	p.Weapons = []*weapon.Weapon{w} // Change the type of p.Weapon to []*weapon.Weapon
	p.WeaponSlots += w.Slots

	return nil
}

// UnequipWeapon unequips the player's weapon.
func UnequipWeapon(p *player.Player) {
	if p.Weapons != nil {
		p.Weapons = nil
		p.WeaponSlots--
	}
}

// EquipItem equips an item to the player if there are available slots.
func EquipItem(p *player.Player, i *item.Item) error {
	// Implement logic to equip an item similar to EquipWeapon
	return nil
}

// UnequipItem unequips an item from the player.
func UnequipItem(p *player.Player) {
	// Implement logic to unequip an item similar to UnequipWeapon
}
