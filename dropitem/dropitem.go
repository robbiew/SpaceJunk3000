package dropitem

import (
	"spacejunk3000/gear"
	"spacejunk3000/weapon"
)

// Item represents any item in the game.
type Item interface {
	// Any methods common to all items can be declared here
	ItemType() string
}

// WeaponWrapper wraps a weapon item.
type WeaponWrapper struct {
	Weapon *weapon.Weapon
	Item   // Embed the Item interface
}

// GearWrapper wraps a gear item.
type GearWrapper struct {
	Gear *gear.Gear
	Item // Embed the Item interface
}
