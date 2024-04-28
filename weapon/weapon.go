package weapon

import (
	"encoding/json"
	"os"
)

// Weapon represents the characteristics of a game weapon.
type Weapon struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	AmmoType     string `json:"ammo_type,omitempty"`
	AmmoCapacity int    `json:"ammo_capacity,omitempty"`
	FireRate     string `json:"fire_rate,omitempty"`
	Jammed       bool   `json:"jammed,omitempty"`
	Slots        int    `json:"slots"`
	Ammo         int    `json:"ammo,omitempty"`
}

// NewWeapon creates a new weapon with the given attributes.
func NewWeapon(name, weaponType string, ammoType string, ammoCapacity int, fireRate string, jammed bool, slots int) *Weapon {
	return &Weapon{Name: name, Type: weaponType, AmmoType: ammoType, AmmoCapacity: ammoCapacity, FireRate: fireRate, Jammed: jammed, Slots: slots}
}

// LoadWeapons loads weapons from a specified JSON file.
func LoadWeapons(filename string) ([]Weapon, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var weapons []Weapon
	err = json.Unmarshal(bytes, &weapons)
	if err != nil {
		return nil, err
	}
	return weapons, nil
}

// SaveWeapons saves a slice of weapons to a specified JSON file.
func SaveWeapons(filename string, weapons []Weapon) error {
	bytes, err := json.Marshal(weapons)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, os.ModePerm)
}
