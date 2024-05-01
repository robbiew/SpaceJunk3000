package weapon

import (
	"encoding/json"
	"fmt"
	"os"
)

// WeaponWrapper wraps a weapon item.
type WeaponWrapper struct {
	Weapon *Weapon
}

// Weapon represents the characteristics of a game weapon.
type Weapon struct {
	Name           string `json:"name"`
	WeaponTypeName string `json:"type"` // Rename the field to avoid conflict
	AmmoType       string `json:"ammo_type,omitempty"`
	AmmoCapacity   int    `json:"ammo_capacity,omitempty"`
	FireRate       string `json:"fire_rate,omitempty"`
	Jammed         bool   `json:"jammed,omitempty"`
	Slots          int    `json:"slots"`
	Ammo           int    `json:"ammo,omitempty"`
}

// NewWeapon creates a new weapon with the given attributes.
func NewWeapon(name, weaponTypeName, ammoType string, ammoCapacity int, fireRate string, jammed bool, slots, ammo int) *Weapon {
	return &Weapon{
		Name:           name,
		WeaponTypeName: weaponTypeName, // Update the field name here
		AmmoType:       ammoType,
		AmmoCapacity:   ammoCapacity,
		FireRate:       fireRate,
		Jammed:         jammed,
		Slots:          slots,
		Ammo:           ammo,
	}
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

// WeaponType returns the type of the weapon.
func (w *Weapon) WeaponType() string {
	return w.WeaponTypeName // Access the field directly
}

func (g *Weapon) String() string {
	return fmt.Sprintf("Weapon: %s, Type: %s", g.Name, g.WeaponTypeName)
}
