package enemy

import (
	"encoding/json"
	"os"
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

// SaveEnemies saves a slice of enemies to a specified JSON file.
func SaveEnemies(filename string, enemies []Enemy) error {
	bytes, err := json.Marshal(enemies)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, os.ModePerm)
}
