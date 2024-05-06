package game

import (
	"fmt"
	"log"
	"math/rand"
	"spacejunk3000/doorutil"
	"spacejunk3000/dropitem"
	"spacejunk3000/enemy"
	"spacejunk3000/gear"
	"spacejunk3000/implant"
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
	CurrentEnemy    enemy.Enemy
	UsedHealthDrone bool // whether the health drone has been used in the current encounter
	Implants        []implant.Implant
	QuitGame        bool
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

func NewGame(playerName string, charType player.CharacterType, weapons []weapon.Weapon, implants []implant.Implant, enemies []enemy.Enemy) (*Game, error) {
	// Initialize the player
	p, err := InitializePlayer(playerName, weapons, implants)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize player: %v", err)
	}

	// Randomly select a location and enemy
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randomEnemyIndex := random.Intn(len(enemies))
	selectedEnemy := enemies[randomEnemyIndex]

	// Create the Game instance
	game := &Game{
		Player:       p,
		Enemies:      enemies,
		Weapons:      weapons,
		CurrentEnemy: selectedEnemy,
		QuitGame:     false,
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

// SelectCharacterType prompts the user to select a character type.
func SelectCharacterType() player.CharacterType {

	// Clear the screen and display the character selection menu
	doorutil.ClearScreen()
	doorutil.DisplayAnsiFile("assets/selectCrew.ans", false)

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
			return player.Marine
		case '3':
			return player.Empath
		case '4':
			return player.Spy
		case '5':
			return player.Scientist
		case '6':
			return player.Smuggler
		default:
			fmt.Println("Invalid choice, please select a valid character type.")
		}
	}
}

func CombatUI(g *Game) {
	// Get the player's character type and convert to string
	charType := fmt.Sprintf("%v", g.Player.Type)
	doorutil.ClearScreen()
	fmt.Println(doorutil.BgMagenta + "                                         " + doorutil.Reset)
	doorutil.PrintAnsiLoc("assets/"+charType+".ans", 1, 2)
}

// Function to present the user with combat options.
func PresentCombatOptions(g *Game) {
	fmt.Println("\r\nEncounter!")
	fmt.Printf("You've encountered an enemy: %s\r\n", g.CurrentEnemy.Name)
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
		if w.WeaponTypeName == "Ranged Weapon" {
			fmt.Println("[S] Shoot - Ranged Weapon")
			fmt.Println("[R] Reload - Ranged Weapon")
			break
		}
	}
}

// Function to handle an encounter.
func HandleEncounter(g *Game) {

	// Print player information
	fmt.Printf("Name: %s\n", g.Player.Name)
	fmt.Printf("Health: %d\n", g.Player.Health)

	// Display the health record using the DisplayHealthRecord method
	// fmt.Printf("%s", g.Player.DisplayHealthRecord())

	// Show available implants
	fmt.Println("\r\nImplants:")
	if g.Player.Implant.Name != "" {
		fmt.Printf("- %s\r\n", g.Player.Implant.Name)
	} else {
		fmt.Println("- No implants equipped")
	}

	// Show available weapons and their ammo
	fmt.Println("\r\nWeapons:")
	for _, w := range g.Player.Weapons {
		// Check if the weapon is of type "Ranged"
		if w.WeaponTypeName == "Ranged Weapon" {
			fmt.Printf("- %s (Ammo: %d)\r\n", w.Name, w.Ammo)
		} else {
			fmt.Printf("- %s\r\n", w.Name)
		}
	}

	// Print player's equipped gear
	fmt.Println("Equipped Gear:")
	if len(g.Gear) == 0 {
		fmt.Println("- None")
	} else {
		for _, g := range g.Gear {
			fmt.Printf("- %s\n", g.Name)
		}
	}

	// Start the game loop
	for {

		// Display the combat UI
		CombatUI(g)
		// Present combat options
		PresentCombatOptions(g)

		// Handle user choice
		HandleCombatChoice(g)

		// Check if the player is dead or chooses to quit
		if g.Player.Health <= 0 {
			fmt.Println("Game Over! You are dead.")
			return
		}

		// Check if the player chooses to quit
		if g.QuitGame {
			// Prompt for playing again
			choice, err := doorutil.PromptYesNo("Quitting will end the game. Are you sure you want to quit?")
			if err != nil {
				log.Println("Error reading keyboard input:", err)
				break
			}

			// Check if the user wants to quit
			if choice == "n" || choice == "N" {
				g.QuitGame = false
				continue
			}
			if choice == "y" || choice == "Y" {
				break
			} else {
				fmt.Println("Invalid choice. Please enter 'y' or 'n'.")
				continue
			}

		}

		// Check if the enemy is dead
		// if g.CurrentEnemy.Health <= 0 {
		//  fmt.Printf("You defeated the %s!\n", g.CurrentEnemy.Name)
		//  return
		// }
	}
}

// HandleCombatChoice handles user's combat choice including selecting an implant if needed.
func HandleCombatChoice(g *Game) {
	g.QuitGame = false
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
			// Iterate over the dropped items and print thssem
			for _, item := range items {
				fmt.Println(item) // Print the dropped item
				fmt.Println("Do you want to pick up this item? (Y/N)")
				choice, _, err := keyboard.GetSingleKey()
				if err != nil {
					fmt.Println("Error reading keyboard input:", err)
					continue // Continue to loop for valid input
				}
				switch choice {
				case 'Y', 'y':
					// Check the underlying type of item
					switch item := item.(type) {
					case *dropitem.WeaponWrapper:
						// Handle weapon
						weapon := item.Weapon
						weaponType := weapon.WeaponType()
						fmt.Printf("Weapon type: %s\n", weaponType)
						// Perform other actions specific to weapons
					case *dropitem.GearWrapper:
						// Handle gear
						gear := item.Gear
						gearType := gear.GearType()
						fmt.Printf("Gear type: %s\n", gearType)
						if err := g.Player.EquipGear(gear); err != nil {
							fmt.Println("Error equipping gear:", err)
							// Handle error (e.g., inform the player)
						} else {
							fmt.Println("Gear equipped successfully:", gear.Name)
						}
					default:
						// Handle unknown item type
						fmt.Println("Unknown item type:", item)
					}
				}
			}

		case 'Q', 'q':
			// Quit the game
			g.QuitGame = true
			return // Exit the function, effectively ending the game loop

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
			ShootWithRangedWeapon(g)
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
			continue // Continue to loop for valid input
		}

		if g.Player.Health <= 0 {
			fmt.Println("Game Over! You are dead.")
			break
		}
		if g.QuitGame {
			break
		}
	}
}

// At the start of each new encounter, you need to reset the UsedHealthDrone field
func StartNewEncounter(g *Game) {
	// Reset the health drone availability for the new encounter
	g.UsedHealthDrone = false

	// Declare quitGame variable
	g.QuitGame = false

	// Continue with encounter setup...
	HandleEncounter(g)
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

// ShootWithRangedWeapon simulates shooting with a ranged weapon.
func ShootWithRangedWeapon(g *Game) {
	// Check if the player has a ranged weapon
	hasRangedWeapon := false
	for _, w := range g.Player.Weapons {
		if w.WeaponTypeName == "Ranged Weapon" {
			hasRangedWeapon = true
			break
		}
	}
	if !hasRangedWeapon {
		fmt.Println("You do not have a ranged weapon.")
		return
	}

	// Select the ranged weapon to use if the player has multiple
	var selectedWeapon *weapon.Weapon
	if len(g.Player.Weapons) > 1 {
		// Implement logic to let the player select a weapon
		// For simplicity, we'll just select the first ranged weapon
		for _, w := range g.Player.Weapons {
			if w.WeaponTypeName == "Ranged Weapon" {
				selectedWeapon = w
				break
			}
		}
	} else {
		selectedWeapon = g.Player.Weapons[0]
	}

	// Check if the player has enough ammo for the required fire rate of the selected weapon
	fmt.Printf("Selected weapon: %s, Ammo: %d, Fire Rate: %d \n", selectedWeapon.Name, selectedWeapon.Ammo, selectedWeapon.FireRate)
	if selectedWeapon.Ammo < selectedWeapon.FireRate {
		fmt.Println("You do not have enough ammo for this weapon.")
		return
	}

	// Select the fire rate
	// For simplicity, we'll assume the fire rate is always 1
	fireRate := 1

	// Fire the weapon and deplete the ammo
	selectedWeapon.Ammo -= fireRate

	// Save the player's updated data after firing
	if err := player.SavePlayer(g.Player); err != nil {
		fmt.Printf("Error saving player data: %v\n", err)
	}

	// Roll ammo dice for each ammo fired
	// For simplicity, we'll just simulate the dice roll without actual dice mechanics
	// We'll use a random number generator to simulate the dice roll
	ammoHits := make(map[string]int)
	for i := 0; i < fireRate; i++ {
		// Simulate dice roll for each ammo type
		// For simplicity, we'll just assume hits for now
		ammoType := selectedWeapon.AmmoType
		hit := simulateHit()
		if hit {
			fmt.Println("Shot hit!")
			ammoHits[ammoType]++
			applyRandomVulnerabilityReduction(&g.Enemies[0], ammoType)
		} else {
			fmt.Println("Shot missed.")
		}
	}

	// Apply ammo hits to enemies
	for _, enemy := range g.Enemies {
		// Check if the enemy is within range
		// For simplicity, we'll assume all enemies are within range
		// Apply damage based on ammo hits and enemy's vulnerabilities
		switch selectedWeapon.AmmoType {
		case "Energy":
			enemy.StrDie -= ammoHits["Energy"] * enemy.EnemyEnerDamage
		case "Ballistic":
			enemy.StrDie -= ammoHits["Ballistic"] * enemy.EnemyBallDamage
		case "Explosive":
			enemy.StrDie -= ammoHits["Explosive"] * enemy.EnemyExplDamage
		}
	}

	// Continue combat logic...
}

// applyRandomVulnerabilityReduction reduces one of the enemy's vulnerabilities corresponding to the player's ammo type.
func applyRandomVulnerabilityReduction(enemy *enemy.Enemy, ammoType string) {
	// Determine which vulnerability to reduce based on the ammo type
	switch ammoType {
	case "Energy":
		if enemy.StrDie > 0 {
			enemy.StrDie--
		}
	case "Ballistic":
		if enemy.DexDie > 0 {
			enemy.DexDie--
		}
	case "Explosive":
		if enemy.IntDie > 0 {
			enemy.IntDie--
		}
	}
}

// simulateHit simulates whether a shot hits or misses.
func simulateHit() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0 // 50% chance of hit or miss
}
