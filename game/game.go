package game

import (
	"fmt"
	"math/rand"
	"spacejunk3000/doorutil"
	"spacejunk3000/enemy"
	"spacejunk3000/item"
	"spacejunk3000/location"
	"spacejunk3000/player"
	"spacejunk3000/weapon"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

type Game struct {
	Player       *player.Player
	Enemies      []enemy.Enemy
	Weapons      []weapon.Weapon
	Location     location.Location
	CurrentEnemy enemy.Enemy
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

func NewGame(playerName string, charType player.CharacterType, weapons []weapon.Weapon, locations []location.Location, enemies []enemy.Enemy) (*Game, error) {
	// Initialize the player
	p, err := InitializePlayer(playerName, weapons)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize player: %v", err)
	}

	// Randomly select a location and enemy
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randomLocationIndex := random.Intn(len(locations))
	randomEnemyIndex := random.Intn(len(enemies))
	selectedLocation := locations[randomLocationIndex]
	selectedEnemy := enemies[randomEnemyIndex]

	// Create the Game instance
	game := &Game{
		Player:       p,
		Enemies:      enemies,
		Weapons:      weapons,
		Location:     selectedLocation,
		CurrentEnemy: selectedEnemy,
	}

	return game, nil
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
	fmt.Println("Choose your character type:")
	fmt.Println("1. Pirate")
	fmt.Println("2. Space Marine")
	fmt.Println("3. Empath")

	for {
		// Initialize keyboard listener
		err := keyboard.Open()
		if err != nil {
			fmt.Println("Error opening keyboard:", err)
			return ""
		}
		defer keyboard.Close()

		// Listen for single key press
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			continue
		}

		// Convert the pressed key to character type
		switch char {
		case '1':
			return player.Pirate
		case '2':
			return player.SpaceMarine
		case '3':
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
	p.Weapons = []*weapon.Weapon{w}
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

// Function to present the user with combat options.
func PresentCombatOptions(g *Game) {
	fmt.Println("\r\nEncounter!")
	fmt.Printf("You've encountered an enemy at %s: %s\r\n", g.Location.Name, g.CurrentEnemy.Name)
	fmt.Println("\r\nChoose your action:")

	fmt.Println("[Q] Quit - run away, lose health, items)")
	fmt.Println("[D] Defend")
	fmt.Println("[U] Use an item or a tech implant (if you have any)")
	fmt.Println("[F] Fight - hand to hand")

	// Check if the player has a ranged weapon
	for _, w := range g.Player.Weapons {
		if w.Type == "Ranged Weapon" {
			fmt.Println("[S] Shoot - Ranged Weapon")
			fmt.Println("[R] Reload - if weapons are out of ammo")
			break
		}
	}
}

// Function to handle user's combat choice.
func HandleCombatChoice(g *Game) {
	for {
		fmt.Println("\r\nChoose your action:")

		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			continue // Continue to loop for valid input
		}

		switch char {
		case ('F' | 'f'):
			// Hand to hand combat logic
			fmt.Println("You chose hand to hand combat.")
		case ('Q' | 'q'):
			// Run away logic
			fmt.Println("You chose to run away.")
		case ('D' | 'd'):
			// Defend logic
			fmt.Println("You chose to defend.")
		case ('R' | 'r'):
			// Reload logic
			fmt.Println("You chose to reload.")
		case ('U' | 'u'):
			// Use item or tech implant logic
			fmt.Println("You chose to use an item or a tech implant.")
		case ('S' | 's'):
			// Ranged combat logic
			fmt.Println("You chose to shoot.")
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
			continue // Continue to loop for valid input
		}

		// If a valid choice is made, break out of the loop
		break
	}
}

// Function to handle an encounter.
func HandleEncounter(g *Game) {
	doorutil.ClearScreen()
	// Print player information
	fmt.Printf("Player Name: %s\n", g.Player.Name)
	fmt.Printf("Health: %d\n", g.Player.Health)

	// Display the health record using the DisplayHealthRecord method
	fmt.Printf("%s", g.Player.DisplayHealthRecord())

	// Show available weapons and their ammo
	fmt.Println("\r\nEquipped Weapons:")
	for _, w := range g.Player.Weapons {
		// Check if the weapon is of type "Ranged"
		if w.Type == "Ranged Weapon" {
			fmt.Printf("- %s (Ammo: %d)\r\n", w.Name, w.Ammo)
		} else {
			fmt.Printf("- %s\r\n", w.Name)
		}
	}

	// Present combat options
	PresentCombatOptions(g)

	// Handle user choice
	HandleCombatChoice(g)
}

// Function to get user's choice.
func GetUserChoice(g *Game) int {
	// Initialize keyboard listener
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Loop until a valid choice is made
	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		// Convert the pressed key to integer if possible
		choice, err := strconv.Atoi(string(char))
		if err == nil && choice >= 1 && choice <= 6 {
			return choice // Return valid choice
		}

		fmt.Println("Invalid choice. Please select a valid option.")
		PresentCombatOptions(g) // Present combat options again
	}
}
