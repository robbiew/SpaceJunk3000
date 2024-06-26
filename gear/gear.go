package gear

import (
	"encoding/json"
	"fmt"
	"os"
)

// Item represents an item in the game.
type Gear struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Slots        int    `json:"slots"`
	GearTypeName string `json:"type"`
	Heal         int    `json:"heal,omitempty"`
	DamageType   string `json:"damage_type,omitempty"`
	SingleUse    bool   `json:"single_use"`
}

// NewItem creates a new item with the given attributes.
func NewGear(name, description string, slots int, gearTypeName string, heal int, damageType string, singleUse bool) *Gear {
	return &Gear{
		Name:         name,
		Description:  description,
		Slots:        slots,
		GearTypeName: gearTypeName,
		Heal:         heal,
		DamageType:   damageType,
		SingleUse:    singleUse,
	}
}

// LoadGear loads items from a specified JSON file.
func LoadGear(filename string) ([]Gear, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var items []Gear
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GearType returns the type of the weapon.
func (g *Gear) GearType() string {
	return g.GearTypeName
}

func (g *Gear) String() string {
	return fmt.Sprintf("Gear: %s, Type: %s", g.Name, g.GearTypeName)
}
