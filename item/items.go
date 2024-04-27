// item.go

package item

import (
	"encoding/json"
	"os"
)

// Item represents an item in the game.
type Item struct {
	Name       string `json:"name"`
	Slots      int    `json:"slots"`
	Type       string `json:"type"`
	Heal       int    `json:"heal,omitempty"`
	DamageType string `json:"damage_type,omitempty"`
}

// NewItem creates a new item with the given attributes.
func NewItem(name string, slots int, itemType string, heal int, damageType string) *Item {
	return &Item{Name: name, Slots: slots, Type: itemType, Heal: heal, DamageType: damageType}
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

// SaveItems saves a slice of items to a specified JSON file.
func SaveItems(filename string, items []Item) error {
	bytes, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, os.ModePerm)
}
