package game

import (
	"time"

	"math/rand"

	"github.com/davidhorak/space-wars/kernel/physics"
)

type Status string

const (
	Initialized Status = "initialized"
	Running     Status = "running"
	Paused      Status = "paused"
	Ended       Status = "ended"
)

type DamageType string

const (
	DamageTypeUnknown DamageType = "unknown"
	DamageTypeLaser   DamageType = "laser"
	DamageTypeRocket  DamageType = "rocket"
)

type Game struct {
	seed             int64
	status           Status
	size             physics.Size
	manager          GameManager
	gracefulEndTimer float64
}

func NewGame(size physics.Size, seed int64) *Game {
	game := Game{
		status:  Initialized,
		size:    size,
		seed:    seed,
		manager: NewGameManager(),
	}

	asteroids := SeedAsteroids(rand.New(rand.NewSource(seed)), game.size.Width, game.size.Height, 1000)
	game.manager.AddGameObjects(asteroids)

	return &game
}

func (game *Game) Status() Status {
	return game.status
}

func (game *Game) Start() {
	if game.status == Running {
		return
	}

	game.status = Running
	game.manager.Logger().GameState(time.Now(), Running)
}

func (game *Game) Pause() {
	if game.status != Running {
		return
	}

	game.status = Paused
	game.manager.Logger().GameState(time.Now(), Paused)
}

func (game *Game) Reset() {
	game.manager.Reset()
	game.manager.Logger().Clear()
}

func (game *Game) Update(deltaTimeMs float64) {
	for _, gameObject := range game.manager.GameObjects() {
		if !gameObject.Enabled() {
			continue
		}
		gameObject.Update(deltaTimeMs, &game.manager)

		// Wrap around the screen
		position := gameObject.Position()
		if position.X < 0 {
			gameObject.SetPosition(physics.Vector2{X: game.size.Width - position.X, Y: position.Y})
		} else if position.X > game.size.Width {
			gameObject.SetPosition(physics.Vector2{X: position.X - game.size.Width, Y: position.Y})
		}
		if position.Y < 0 {
			gameObject.SetPosition(physics.Vector2{X: position.X, Y: game.size.Height - position.Y})
		} else if position.Y > game.size.Height {
			gameObject.SetPosition(physics.Vector2{X: position.X, Y: position.Y - game.size.Height})
		}
	}

	size := game.manager.GameObjectSize()
	for i := 0; i < size-1; i++ {
		a := game.manager.GetGameObjectByIndex(i)
		colliderA := a.Collider()
		if !a.Enabled() || colliderA == nil {
			continue
		}
		for j := i + 1; j < size; j++ {
			b := game.manager.GetGameObjectByIndex(j)
			colliderB := b.Collider()
			if !b.Enabled() || colliderB == nil {
				continue
			}

			if colliderA.CollidesWith(colliderB) {
				a.OnCollision(b, &game.manager, 0)
				b.OnCollision(a, &game.manager, 1)
			}
		}
	}

	if game.manager.HasEnded(deltaTimeMs) {
		game.status = Ended
		game.manager.Logger().GameState(time.Now(), Ended)
	}
}

func (game *Game) SpaceshipAction(name string, action func(spaceShip *Spaceship, gameManager *GameManager)) error {
	spaceShip, err := game.manager.GetSpaceship(name)
	if err != nil {
		return err
	}
	action(spaceShip, &game.manager)
	return nil
}

func (game *Game) AddSpaceship(name string, position physics.Vector2, rotation float64) error {
	spaceShip := NewSpaceship(NewUUID(), name, position, rotation)
	return game.manager.AddSpaceship(spaceShip)
}

func (game *Game) RemoveSpaceship(name string) error {
	return game.manager.RemoveSpaceship(name)
}

func (game *Game) Serialize() map[string]interface{} {
	gameObjects := make([]interface{}, 0)
	for _, gameObject := range game.manager.GameObjects() {
		gameObjects = append(gameObjects, gameObject.Serialize())
	}

	logs := make([]interface{}, 0)
	for _, log := range game.manager.Logger().Logs() {
		logs = append(logs, log.Serialize())
	}

	return map[string]interface{}{
		"status": string(game.Status()),
		"seed":   game.seed,
		"size": map[string]interface{}{
			"width":  game.size.Width,
			"height": game.size.Height,
		},
		"gameObjects": gameObjects,
		"logs":        logs,
	}
}
