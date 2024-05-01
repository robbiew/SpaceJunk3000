package enemy

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"spacejunk3000/gear"
	"spacejunk3000/weapon"
)

// Enemy represents the characteristics of a game enemy.
type Enemy struct {
	Name               string `json:"name"`
	Health             int    `json:"health"`
	Damage             int    `json:"damage"`
	Desc               string `json:"desc"`
	ToHitMightDie      int    `json:"toHitMightDie"`
	ToHitCunningDie    int    `json:"toHitCunningDie"`
	ToHitWisdomDie     int    `json:"toHitWisdomDie"`
	EnemyBallDamage    int    `json:"enemyBallDamage"`
	EnemyEnerDamage    int    `json:"enemyEnerDamage"`
	EnemyExplDamage    int    `json:"enemyExplDamage"`
	PlayerRangedDamage int    `json:"playerRangedDamage"`
	PlayerCloseDamage  int    `json:"playerCloseDamage"`
	ItemDrop           int    `json:"itemDrop"`
	Initiative         bool   `json:"initiative"`
}

type Item interface {
	// Define any common methods or fields here
	String() string // Define a method common to both weapons and gear
}

// NewEnemy creates a new enemy with the given attributes.
func NewEnemy(name string, health, damage int, desc string, toHitMightDie, toHitCunningDie, toHitWisdomDie, enemyBallDamage, enemyEnerDamage, enemyExplDamage, playerRangedDamage, playerCloseDamage, itemDrop int, initiative bool) *Enemy {
	return &Enemy{
		Name:               name,
		Health:             health,
		Damage:             damage,
		Desc:               desc,
		ToHitMightDie:      toHitMightDie,
		ToHitCunningDie:    toHitCunningDie,
		ToHitWisdomDie:     toHitWisdomDie,
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

func (e *Enemy) DropItems() ([]Item, error) {
	items := make([]Item, 0)

	// Determine the number of items to drop based on ItemDrop field
	for i := 0; i < e.ItemDrop; i++ {
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
				items = append(items, WeaponWrapper{Weapon: weapons[randomIndex]})
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
				items = append(items, GearWrapper{Gear: gears[randomIndex]})
			}
		}
	}

	return items, nil
}

// WeaponWrapper wraps a weapon for additional functionality.
type WeaponWrapper struct {
	weapon.Weapon
}

// String returns the string representation of the weapon.
func (w WeaponWrapper) String() string {
	return fmt.Sprintf("Weapon: %s, Type: %s", w.Name, w.WeaponTypeName)
}

// GearWrapper wraps a gear for additional functionality.
type GearWrapper struct {
	gear.Gear
}

// String returns the string representation of the gear.
func (g GearWrapper) String() string {
	return fmt.Sprintf("Gear: %s, Type: %s", g.Name, g.GearTypeName)
}
