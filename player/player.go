package player

import (
	"encoding/json"
	"fmt"
	"os"
	"spacejunk3000/door"
	"spacejunk3000/gear"
	"spacejunk3000/implant"
	"spacejunk3000/weapon"
)

// Define the Player struct with exported Inventory field
type Player struct {
	Name         string           `json:"name"`   // Exported field
	Type         CharacterType    `json:"type"`   // Exported field
	Health       int              `json:"health"` // Exported field
	HealthRecord []string         `json:"health_record"`
	Stats        Stats            `json:"stats"`            // Exported field
	Alive        bool             `json:"alive"`            // Unexported field
	TimeLeft     int              `json:"-"`                // Unexported field
	Emulation    int              `json:"-"`                // Unexported field
	NodeNum      int              `json:"-"`                // Unexported field
	Weapons      []*weapon.Weapon `json:"weapon,omitempty"` // Include a field for the weapon
	WeaponSlots  int              `json:"weapon_slots"`     // Number of filled weapon slots
	Gear         []*gear.Gear     `json:"gear"`             // Include a field for the gear
	GearSlots    int              `json:"gear_slots"`       // Number of filled item slots
	MaxSlots     int              `json:"max_slots"`        // Maximum number of total slots
	CrewDice     CrewDice         `json:"crew_dice"`
	Implant      implant.Implant  `json:"implant"` // Include a field for the implants

}

type CharacterType string

const (
	Pirate    CharacterType = "Pirate"    // Miller
	Marine    CharacterType = "Marine"    // Smith
	Empath    CharacterType = "Empath"    // Tailor
	Spy       CharacterType = "Spy"       // Tanner
	Scientist CharacterType = "Scientist" // Abbot
	Smuggler  CharacterType = "Smuggler"  // Cook

)

type Stats struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Intelligence int `json:"intelligence"`
}

// CrewDice represents the crew dice associated with each character type.
type CrewDice struct {
	DieSide1 string `json:"die_side_1"`
	DieSide2 string `json:"die_side_2"`
	DieSide3 string `json:"die_side_3"`
	DieSide4 string `json:"die_side_4"`
	DieSide5 string `json:"die_side_5"`
	DieSide6 string `json:"die_side_6"`
}

func PrintPlayerInventory(player *Player) {
	// Create a slice to hold all items (weapons and gear)
	items := make([]interface{}, 0, player.MaxSlots)

	// Append weapons to the items slice
	for _, weapon := range player.Weapons {
		items = append(items, weapon)
	}

	// Append gear to the items slice
	for _, gear := range player.Gear {
		items = append(items, gear)
	}

	// Print items
	for i := 0; i < player.MaxSlots; i++ {
		if i < len(items) {
			// Print actual item
			switch item := items[i].(type) {
			case *weapon.Weapon:
				fmt.Printf("%s%d %s%-14s %s%-2d %s%-9s %s%-4d %-4d\r\n", door.BlackHi, i+1, door.CyanHi, item.Name, door.Reset, item.Slots, door.Cyan, item.WeaponTypeName, door.Reset, item.Ammo, item.FireRate)
			case *gear.Gear:
				fmt.Printf("%d %-14s %-2d %-9s %-4s %-4s\r\n", i+1, item.Name, item.Slots, item.GearTypeName, "-", "-")
			}
		} else {
			// Print empty row
			fmt.Printf("%s%d %-14s %-2s %-9s %-4s %-4s%s\r\n", door.BlackHi, i+1, "-", "-", "-", "-", "-", door.Reset)
		}
	}

}

// NewPlayer creates a new player instance with the provided attributes.
func NewPlayer(name string, charType CharacterType, timeLeft int, nodeNum int, emulation int) (*Player, error) {
	// Get character stats based on character type
	stats, err := GetCharacterStats(charType)
	if err != nil {
		return nil, fmt.Errorf("failed to get character stats: %v", err)
	}
	dice, err := GetCrewDice(charType)
	if err != nil {
		return nil, fmt.Errorf("failed to get crew dice: %v", err)
	}
	// Initialize the health record with all "-" for full health
	healthRecord := make([]string, 12)
	for i := range healthRecord {
		healthRecord[i] = "-"
	}

	return &Player{
		Name:         name,
		Type:         charType,
		Health:       12,
		HealthRecord: healthRecord,
		Stats:        stats,
		TimeLeft:     timeLeft,
		NodeNum:      nodeNum,
		CrewDice:     dice,
		Emulation:    emulation,
		Alive:        true,
		MaxSlots:     4,                         // Default value, can be modified if needed
		Weapons:      make([]*weapon.Weapon, 0), // Initialize the weapons slice
		WeaponSlots:  0,                         // Initialize the weapon slots
		Implant:      implant.Implant{},         // Initialize the implant
		Gear:         make([]*gear.Gear, 0),     // Initialize the gear slice
		GearSlots:    0,                         // Initialize the gear slots

	}, nil
}

// GetCharacterStats returns the stats associated with the provided character type.
func GetCharacterStats(charType CharacterType) (Stats, error) {
	switch charType {
	case Pirate: // Miller
		return Stats{Strength: 3, Dexterity: 4, Intelligence: 1}, nil
	case Marine: // Smith
		return Stats{Strength: 4, Dexterity: 1, Intelligence: 3}, nil
	case Empath: // Tailor
		return Stats{Strength: 1, Dexterity: 4, Intelligence: 3}, nil
	case Spy: // Tanner
		return Stats{Strength: 1, Dexterity: 3, Intelligence: 4}, nil
	case Scientist: // Abbot
		return Stats{Strength: 3, Dexterity: 1, Intelligence: 4}, nil
	case Smuggler: // Cook
		return Stats{Strength: 4, Dexterity: 3, Intelligence: 1}, nil
	default:
		return Stats{}, fmt.Errorf("try again")
	}
}

// getCrewDice returns the CrewDice associated with the provided character type.
func GetCrewDice(charType CharacterType) (CrewDice, error) {
	switch charType {
	case Pirate: // Miller
		return CrewDice{
			DieSide1: "dexerity",
			DieSide2: "strength",
			DieSide3: "double dexerity",
			DieSide4: "dexerity",
			DieSide5: "intelligence",
			DieSide6: "ddouble strength",
		}, nil
	case Marine: // Smith
		return CrewDice{
			DieSide1: "strength",
			DieSide2: "dexterity",
			DieSide3: "double intelligence",
			DieSide4: "strength",
			DieSide5: "intelligence",
			DieSide6: "double strength",
		}, nil
	case Empath: // Tailor
		return CrewDice{
			DieSide1: "dexterity",
			DieSide2: "strength",
			DieSide3: "double intelligence",
			DieSide4: "dexterity",
			DieSide5: "strength",
			DieSide6: "double dexterity",
		}, nil
	case Spy: // Tanner
		return CrewDice{
			DieSide1: "intelligence",
			DieSide2: "dexterity",
			DieSide3: "double intelligence",
			DieSide4: "intelligence",
			DieSide5: "dexterity",
			DieSide6: "double dexterity",
		}, nil
	case Scientist: // Abbot
		return CrewDice{
			DieSide1: "intelligence",
			DieSide2: "dexterity",
			DieSide3: "double intelligence",
			DieSide4: "intelligence",
			DieSide5: "dexterity",
			DieSide6: "double strength",
		}, nil
	case Smuggler: // Cook
		return CrewDice{
			DieSide1: "strength",
			DieSide2: "dexterity",
			DieSide3: "double dexterity",
			DieSide4: "strength",
			DieSide5: "dexterity",
			DieSide6: "double strength",
		}, nil
	default:
		return CrewDice{}, fmt.Errorf("unsupported character type")
	}
}

// SavePlayer serializes the player data to JSON and writes it to a file.
func SavePlayer(p *Player) error {
	// Marshal player data to JSON
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error marshaling player data: %v", err)
	}

	// Print out the serialized player data for debugging
	// fmt.Println("Serialized player data:", string(data))

	// Filename based on player name, which is the unique ID
	filename := fmt.Sprintf("data/u-%s.json", p.Name)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing player data to file: %v", err)
	}

	return nil
}

// LoadPlayer deserializes player data from a JSON file.
func LoadPlayer(name string) (*Player, error) {
	filename := fmt.Sprintf("data/u-%s.json", name)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading player data file: %v", err) // File not found could mean new player
	}
	var p Player
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("error unmarshaling player data: %v", err)
	}
	return &p, nil
}

func ResetPlayer(p *Player) {
	p.Health = 12
	p.Alive = true
	p.Weapons = nil
	p.WeaponSlots = 0
	p.GearSlots = 0
	p.MaxSlots = 4

	// Additional reset logic as needed
}

// EquipGear equips a gear to the player if there are available slots.
func (p *Player) EquipGear(g *gear.Gear) error {
	// Check if there are enough gear slots to equip the gear
	if p.GearSlots+g.Slots > p.MaxSlots {
		return fmt.Errorf("cannot carry that much")
	}

	// Equip the gear to the player
	p.Gear = append(p.Gear, g)
	p.GearSlots += g.Slots

	// Save the player's data after equipping the weapon
	if err := SavePlayer(p); err != nil {
		return fmt.Errorf("failed to save player: %v", err)
	}

	return nil
}

// EquipWeapon equips a weapon to the player if there are available slots.
func (p *Player) EquipWeapon(w *weapon.Weapon) error {
	// Check if there are enough weapon slots to equip the weapon
	if p.WeaponSlots+w.Slots > p.MaxSlots {
		return fmt.Errorf("cannot carry that much")
	}

	// Equip the weapon to the player
	p.Weapons = append(p.Weapons, w)
	p.WeaponSlots += w.Slots

	// Save the player's data after equipping the weapon
	if err := SavePlayer(p); err != nil {
		return fmt.Errorf("failed to save player: %v", err)
	}

	return nil
}

// UnequipWeapon unequips the player's weapon.
func (p *Player) UnequipWeapon() {
	if len(p.Weapons) > 0 {
		p.Weapons = p.Weapons[:len(p.Weapons)-1]
		p.WeaponSlots--
	}
	// Save the player's data after equipping the weapon
	if err := SavePlayer(p); err != nil {
		fmt.Printf("failed to save player: %v", err)
	}

}

// AdjustHealth updates the player's health and modifies the health record.
// Pass a positive number to heal, or a negative number to deal damage.
func (p *Player) AdjustHealth(amount int) {
	for i := 0; i < len(p.HealthRecord) && amount != 0; i++ {
		if amount > 0 && p.HealthRecord[i] == "\\" {
			p.HealthRecord[i] = "/"
			amount--
		} else if amount < 0 && p.HealthRecord[i] == "-" {
			p.HealthRecord[i] = "\\"
			amount++
		}
	}
	// Update the actual Health value accordingly
	p.Health += amount
}

// DisplayHealthRecord outputs the player's health record as a string.
func (p *Player) DisplayHealthRecord() string {
	// Top part of the medical record
	record := "\r\n┌──────────────────────────────┐\r\n│ MEDICAL RECORD               │\n"

	// Adding the health points with appropriate symbols
	for i := 12; i > 0; i-- {
		line := fmt.Sprintf("│ %2d │ ", i)
		for j := 0; j < 12; j++ {
			if j < p.Health {
				// Check the HealthRecord for visual representation of health points
				line += p.HealthRecord[j] + " "
			} else {
				line += "  " // Empty space for lost health points
			}
		}
		record += line + "│\r\n"
	}

	// Bottom part of the medical record
	record += "└──────────────────────────────┘"

	return record
}

// UpdateHealth updates the player's health and handles death.
func UpdateHealth(p *Player, delta int) {
	p.Health += delta
	if p.Health <= 0 {
		p.Alive = false
		SavePlayer(p) // Save the dead state
		fmt.Println("You have died and must start over.")
		ResetPlayer(p) // Reset for new game start
	}
	SavePlayer(p) // Save any changes to player data
}
