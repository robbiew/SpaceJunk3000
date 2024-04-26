package player

import "fmt"

type CharacterType string

const (
	Pirate      CharacterType = "Pirate"
	SpaceMarine CharacterType = "Space Marine"
	Empath      CharacterType = "Empath"
)

type Stats struct {
	Might   int
	Cunning int
	Wisdom  int
}

type Player struct {
	Name      string
	Type      CharacterType
	Health    int
	Stats     Stats
	Inventory []string // Simple inventory system for now
}

// NewPlayer initializes a new player with default values and handles potential errors.
func NewPlayer(name string, charType CharacterType) (*Player, error) {
	if name == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}
	// Optionally, add more validation if necessary, for example:
	// if charType is not one of the predefined ones, return an error.

	stats := Stats{Might: 4, Cunning: 4, Wisdom: 4}
	return &Player{Name: name, Type: charType, Health: 12, Stats: stats}, nil
}
