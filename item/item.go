// item.go

package item

import (
	"encoding/json"
	"os"
)

// Item represents an item in the game.
type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Slots       int    `json:"slots"`
	Type        string `json:"type"`
	Heal        int    `json:"heal,omitempty"`
	DamageType  string `json:"damage_type,omitempty"`
	SingleUse   bool   `json:"single_use"`
}

// NewItem creates a new item with the given attributes.
func NewItem(name, description string, slots int, itemType string, heal int, damageType string, singleUse bool) *Item {
	return &Item{
		Name:        name,
		Description: description,
		Slots:       slots,
		Type:        itemType,
		Heal:        heal,
		DamageType:  damageType,
		SingleUse:   singleUse,
	}
}

// LoadItems loads items from a specified JSON file.
func LoadItems(filename string) ([]Item, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var items []Item
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
