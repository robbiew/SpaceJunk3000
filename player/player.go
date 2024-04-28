package player

import (
	"encoding/json"
	"fmt"
	"os"
	"spacejunk3000/weapon"
)

// Define the Player struct with exported Inventory field
type Player struct {
	Name        string           `json:"name"`             // Exported field
	Type        CharacterType    `json:"type"`             // Exported field
	Health      int              `json:"health"`           // Exported field
	Stats       Stats            `json:"stats"`            // Exported field
	Alive       bool             `json:"alive"`            // Unexported field
	TimeLeft    int              `json:"-"`                // Unexported field
	Emulation   int              `json:"-"`                // Unexported field
	NodeNum     int              `json:"-"`                // Unexported field
	Weapons     []*weapon.Weapon `json:"weapon,omitempty"` // Include a field for the weapon
	WeaponSlots int              `json:"weapon_slots"`     // Number of filled weapon slots
	ItemSlots   int              `json:"item_slots"`       // Number of filled item slots
	MaxSlots    int              `json:"max_slots"`        // Maximum number of total slots
	CrewDice    CrewDice         `json:"crew_dice"`
}

type CharacterType string

const (
	Pirate      CharacterType = "Pirate"
	SpaceMarine CharacterType = "Space Marine"
	Empath      CharacterType = "Empath"
)

type Stats struct {
	Might   int `json:"might"`
	Cunning int `json:"cunning"`
	Wisdom  int `json:"wisdom"`
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

	return &Player{
		Name:      name,
		Type:      charType,
		Health:    12,
		Stats:     stats,
		TimeLeft:  timeLeft,
		NodeNum:   nodeNum,
		CrewDice:  dice,
		Emulation: emulation,
		Alive:     true,
		MaxSlots:  4,                         // Default value, can be modified if needed
		Weapons:   make([]*weapon.Weapon, 0), // Initialize the weapons slice

	}, nil
}

// GetCharacterStats returns the stats associated with the provided character type.
func GetCharacterStats(charType CharacterType) (Stats, error) {
	switch charType {
	case Pirate:
		return Stats{Might: 1, Cunning: 1, Wisdom: 1}, nil
	case SpaceMarine:
		return Stats{Might: 2, Cunning: 2, Wisdom: 2}, nil
	case Empath:
		return Stats{Might: 3, Cunning: 3, Wisdom: 3}, nil
	default:
		return Stats{}, fmt.Errorf("unsupported character type")
	}
}

// getCrewDice returns the CrewDice associated with the provided character type.
func GetCrewDice(charType CharacterType) (CrewDice, error) {
	switch charType {
	case Pirate:
		return CrewDice{
			DieSide1: "might",
			DieSide2: "wisdom",
			DieSide3: "might",
			DieSide4: "double cunning",
			DieSide5: "double might",
			DieSide6: "cunning",
		}, nil
	case SpaceMarine:
		return CrewDice{
			DieSide1: "cunning",
			DieSide2: "wisdom",
			DieSide3: "cunning",
			DieSide4: "double cunning",
			DieSide5: "double might",
			DieSide6: "might",
		}, nil
	case Empath:
		return CrewDice{
			DieSide1: "wisdom",
			DieSide2: "might",
			DieSide3: "wisdom",
			DieSide4: "double cunning",
			DieSide5: "double wisdom",
			DieSide6: "cunning",
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
	p.ItemSlots = 0
	p.MaxSlots = 4

	// Additional reset logic as needed
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

	return nil
}

// UnequipWeapon unequips the player's weapon.
func (p *Player) UnequipWeapon() {
	if len(p.Weapons) > 0 {
		p.Weapons = p.Weapons[:len(p.Weapons)-1]
		p.WeaponSlots--
	}
}
