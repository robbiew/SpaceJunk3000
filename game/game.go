package game

import (
	"fmt"
	"math/rand"
	"spacejunk3000/enemy"
	"spacejunk3000/player"
	"time"
)

type Game struct {
	Player  *player.Player
	Enemies []enemy.Enemy
}

func NewGame(playerName string, charType player.CharacterType) (*Game, error) {
	p, err := player.NewPlayer(playerName, charType)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	enemies, err := enemy.LoadEnemies("data/enemies.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load enemies: %v", err)
	}

	return &Game{Player: p, Enemies: enemies}, nil
}

func (g *Game) Start() error {
	fmt.Println("Game has started.")
	fmt.Printf("Player %s has entered the game as a %s with %d health points.\n", g.Player.Name, g.Player.Type, g.Player.Health)

	if len(g.Enemies) == 0 {
		return fmt.Errorf("no enemies available for an encounter")
	}

	// Create a new local random generator
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Use the local random generator
	randomIndex := r.Intn(len(g.Enemies))
	randomEnemy := g.Enemies[randomIndex]

	fmt.Printf("A wild %s appears! It has %d health and can deal %d damage.\n", randomEnemy.Name, randomEnemy.Health, randomEnemy.Damage)

	// Implement encounter logic
	// if err := g.handleEncounter(randomEnemy); err != nil {
	//     return fmt.Errorf("encounter failed: %v", err)
	// }

	return nil
}
