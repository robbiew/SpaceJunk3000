package enemy

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"spacejunk3000/dropitem"
	"spacejunk3000/gear"
	"spacejunk3000/weapon"
)

// Enemy represents the characteristics of a game enemy.
type Enemy struct {
	Name               string `json:"name"`
	Desc               string `json:"desc"`
	StrDie             int    `json:"strDie"`
	DexDie             int    `json:"dexDie"`
	IntDie             int    `json:"IntDie"`
	EnemyBallDamage    int    `json:"enemyBallDamage"`
	EnemyEnerDamage    int    `json:"enemyEnerDamage"`
	EnemyExplDamage    int    `json:"enemyExplDamage"`
	PlayerRangedDamage int    `json:"playerRangedDamage"`
	PlayerCloseDamage  int    `json:"playerCloseDamage"`
	ItemDrop           int    `json:"itemDrop"`
	Initiative         bool   `json:"initiative"`
}

// NewEnemy creates a new enemy with the given attributes.
func NewEnemy(name string, health, damage int, desc string, strDie, dexDie, intDie, enemyBallDamage, enemyEnerDamage, enemyExplDamage, playerRangedDamage, playerCloseDamage, itemDrop int, initiative bool) *Enemy {
	return &Enemy{
		Name:               name,
		Desc:               desc,
		StrDie:             strDie,
		DexDie:             dexDie,
		IntDie:             intDie,
		EnemyBallDamage:    enemyBallDamage,
		EnemyEnerDamage:    enemyEnerDamage,
		EnemyExplDamage:    enemyExplDamage,
		PlayerRangedDamage: playerRangedDamage,
		PlayerCloseDamage:  playerCloseDamage,
		ItemDrop:           itemDrop,
		Initiative:         initiative,
	}
}

// LoadEnemies loads enemies from a specified JSON file.
func LoadEnemies(filename string) ([]Enemy, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var enemies []Enemy
	err = json.Unmarshal(bytes, &enemies)
	if err != nil {
		return nil, err
	}
	return enemies, nil
}

// DropItems returns a single item dropped by the enemy.
func (e *Enemy) DropItems() ([]dropitem.Item, error) {
	// Randomly choose between dropping a weapon or gear
	if rand.Intn(2) == 0 {
		// Enemy drops a weapon
		weapons, err := weapon.LoadWeapons("data/weapons.json")
		if err != nil {
			return nil, err
		}
		if len(weapons) > 0 {
			randomIndex := rand.Intn(len(weapons))
			fmt.Println("Adding weapon:", weapons[randomIndex].Name) // Debug print
			return []dropitem.Item{&dropitem.WeaponWrapper{Weapon: &weapons[randomIndex]}}, nil
		}
	} else {
		// Enemy drops gear
		gears, err := gear.LoadGear("data/gear.json")
		if err != nil {
			return nil, err
		}
		if len(gears) > 0 {
			randomIndex := rand.Intn(len(gears))
			fmt.Println("Adding gear:", gears[randomIndex].Name) // Debug print
			return []dropitem.Item{&dropitem.GearWrapper{Gear: &gears[randomIndex]}}, nil
		}
	}

	return nil, nil // No items to drop
}
