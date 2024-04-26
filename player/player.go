package player

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define the Player struct with exported Inventory field
type Player struct {
	Name      string        `json:"name"`      // Exported field
	Type      CharacterType `json:"type"`      // Exported field
	Health    int           `json:"health"`    // Exported field
	Stats     Stats         `json:"stats"`     // Exported field
	Inventory []ItemType    `json:"inventory"` // Exported field
	Alive     bool          `json:"alive"`     // Unexported field
	TimeLeft  int           `json:"-"`         // Unexported field
	Emulation int           `json:"-"`         // Unexported field
	NodeNum   int           `json:"-"`         // Unexported field
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

// Define ItemType as an exported type
type ItemType string

const (
	Sword  ItemType = "Sword"
	Shield ItemType = "Shield"
	// Add more item types as needed
)

// NewPlayer initializes a new player with default values and handles potential errors.
func NewPlayer(name string, charType CharacterType, dropTimeLeft, dropEmulation, nodeNum int) (*Player, error) {
	if name == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}
	// Optionally, add more validation if necessary:
	// if charType is not one of the predefined ones, return an error.

	stats := Stats{Might: 1, Cunning: 2, Wisdom: 4}
	inventory := []ItemType{Sword, Shield}
	return &Player{
		Name:      name,
		Type:      charType,
		Health:    12,
		Stats:     stats,
		Inventory: inventory,
		Alive:     true,
		TimeLeft:  dropTimeLeft,
		Emulation: dropEmulation,
		NodeNum:   nodeNum,
	}, nil
}

// SavePlayer serializes the player data to JSON and writes it to a file.
func SavePlayer(p *Player) error {
	// Marshal player data to JSON
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error marshaling player data: %v", err)
	}

	// Print out the serialized player data for debugging
	fmt.Println("Serialized player data:", string(data))

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
	p.Inventory = []ItemType{} // Reset inventory to empty slice
	// Additional reset logic as needed
}
