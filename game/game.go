package game

import (
	"fmt"
	"math/rand"
	"spacejunk3000/doorutil"
	"spacejunk3000/enemy"
	"spacejunk3000/implant"

	"spacejunk3000/gear"
	"spacejunk3000/location"
	"spacejunk3000/player"
	"spacejunk3000/weapon"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

type Game struct {
	Player          *player.Player
	Enemies         []enemy.Enemy
	Weapons         []weapon.Weapon
	Gear            []gear.Gear
	Location        location.Location
	CurrentEnemy    enemy.Enemy
	UsedHealthDrone bool // whether the health drone has been used in the current encounter
	Implants        []implant.Implant
}

// InitializePlayer initializes a player by loading an existing one or creating a new one if not found.
func InitializePlayer(playerName string, weapons []weapon.Weapon, implants []implant.Implant) (*player.Player, error) {
	// Load existing player or create a new one if not found
	p, err := player.LoadPlayer(playerName)
	if err != nil || p == nil {
		charType := SelectCharacterType()                  // Let the user select a character type if creating a new player
		selectedImplant := implant.SelectImplant(implants) // Select an implant

		// Initialize the player with default values and selected implant
		p, err = player.NewPlayer(playerName, charType, 0, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to create new player: %v", err)
		}

		// Set the selected implant for the player
		p.Implant = selectedImplant

		// Randomly select a weapon for the player
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		randomIndex := random.Intn(len(weapons))

		// Equip the randomly selected weapon to the player
		if err := p.EquipWeapon(&weapons[randomIndex]); err != nil {
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
		if err := p.EquipWeapon(&weapons[randomIndex]); err != nil {
			return nil, fmt.Errorf("failed to equip weapon: %v", err)
		}

		// Save the player data
		if err := player.SavePlayer(p); err != nil {
			return nil, fmt.Errorf("failed to save player: %v", err)
		}
	}

	return p, nil
}

func NewGame(playerName string, charType player.CharacterType, weapons []weapon.Weapon, implants []implant.Implant, locations []location.Location, enemies []enemy.Enemy) (*Game, error) {
	// Initialize the player
	p, err := InitializePlayer(playerName, weapons, implants)
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
func StartGame(playerName string, weapons []weapon.Weapon, implants []implant.Implant) (*player.Player, error) {
	// Initialize the player
	p, err := InitializePlayer(playerName, weapons, implants)
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

// Function to present the user with combat options.
func PresentCombatOptions(g *Game) {
	fmt.Println("\r\nEncounter!")
	fmt.Printf("You've encountered an enemy at %s: %s\r\n", g.Location.Name, g.CurrentEnemy.Name)
	fmt.Println("\r\nChoose your action:")

	fmt.Println("[Q] Quit - run away, lose health, items)")
	fmt.Println("[G] Use Gear")
	fmt.Println("[C] Use Cyber implant")
	if !g.UsedHealthDrone {
		fmt.Println("[H] Activate Health Drone - heal (once per day)")
	} else {
		fmt.Println("[-] Health Drone is unavailable")
	}
	fmt.Println("[F] Fight - hand to hand")

	// Check if the player has a ranged weapon
	for _, w := range g.Player.Weapons {
		if w.Type == "Ranged Weapon" {
			fmt.Println("[S] Shoot - Ranged Weapon")
			fmt.Println("[R] Reload - Ranged Weapon")
			break
		}
	}
}

// HandleCombatChoice handles user's combat choice including selecting an implant if needed.
func HandleCombatChoice(g *Game) {
	for {
		fmt.Println("\r\nChoose your action:")

		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			continue // Continue to loop for valid input
		}

		switch char {
		case 'F', 'f':
			// Hand to hand combat logic
			fmt.Printf("You chose hand to hand combat with %s\n", g.CurrentEnemy.Name)
			items, err := g.CurrentEnemy.DropItems()
			if err != nil {
				// Handle error
				fmt.Println("Error dropping items:", err)
				return
			}
			fmt.Printf("Dropped %d items:\n", len(items)) // Print the number of dropped items
			// Iterate over the dropped items and print them
			for _, item := range items {
				switch v := item.(type) {
				case weapon.Weapon:
					fmt.Println(v.String()) // Print the dropped weapon using its String method
				case gear.Gear:
					fmt.Println(v.String()) // Print the dropped gear using its String method
				default:
					fmt.Println("Unknown item type:", v)
				}
			}

		case 'Q', 'q':
			// Run away logic
			fmt.Println("You chose to run away.")
		case 'G', 'g':
			// Gear logic
			fmt.Println("You chose to use gear.")
		case 'R', 'r':
			// Reload logic
			fmt.Println("You chose to reload.")
		case 'C', 'c':
			// Use implant logic
			fmt.Println("You chose to use an implant.")
			// Check if the player has an implant
			if g.Player.Implant.Name != "" {
				// If the player has an implant, perform actions with it
				fmt.Printf("You selected %s implant.\n", g.Player.Implant.Name)
				// Perform actions with the selected implant if needed
			} else {
				fmt.Println("You don't have any implants.")
			}
		case 'H', 'h':
			if !g.UsedHealthDrone {
				// Activate Health Drone logic here
				fmt.Println("Activating Health Drone.")
				// Update player health here
				g.UsedHealthDrone = true // Mark the drone as used
			} else {
				fmt.Println("Health Drone is unavailable.")
			}
		case 'S', 's':
			// Ranged combat logic
			fmt.Println("You chose to shoot.")
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
			continue // Continue to loop for valid input
		}
		// Break out of the loop if a valid choice is made
		if g.UsedHealthDrone || char == 'F' || char == 'f' || char == 'Q' || char == 'q' || char == 'D' || char == 'd' || char == 'R' || char == 'r' || char == 'U' || char == 'u' || char == 'S' || char == 's' {
			break
		} else {
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}

// At the start of each new encounter, you need to reset the UsedHealthDrone field
func StartNewEncounter(g *Game) {
	// Reset the health drone availability for the new encounter
	g.UsedHealthDrone = true
	// Continue with encounter setup...
	HandleEncounter(g)
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
