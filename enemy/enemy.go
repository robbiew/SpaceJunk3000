package enemy

import (
	"encoding/json"
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

// DropItems randomly generates items dropped by the enemy.
func (e *Enemy) DropItems() ([]interface{}, error) {
	items := make([]interface{}, 0)

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
				// fmt.Println("Adding weapon:", weapons[randomIndex].Name) // Debug print
				items = append(items, weapons[randomIndex])
			}
		} else {
			// Enemy drops gear
			gears, err := gear.LoadGear("data/gear.json")
			if err != nil {
				return nil, err
			}
			if len(gears) > 0 {
				randomIndex := rand.Intn(len(gears))
				// fmt.Println("Adding gear:", gears[randomIndex].Name) // Debug print
				items = append(items, gears[randomIndex])
			}
		}
	}

	return items, nil
}
