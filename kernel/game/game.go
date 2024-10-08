package game

import (
	"fmt"
	"time"

	"encoding/json"
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
	return &Game{
		status:  Initialized,
		size:    size,
		seed:    seed,
		manager: NewGameManager(),
	}
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

func (game *Game) SeedAsteroids() {
	asteroids := SeedAsteroids(rand.New(rand.NewSource(game.seed)), game.size.Width, game.size.Height, 1000)
	game.manager.AddGameObjects(asteroids)
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

func Deserialize(jsonData string) (*Game, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}

	size := data["size"].(map[string]interface{})

	game := NewGame(
		physics.Size{
			Width:  size["width"].(float64),
			Height: size["height"].(float64),
		},
		int64(data["seed"].(float64)),
	)

	uuid := int64(0)

	// Game objects
	for _, gameObject := range data["gameObjects"].([]interface{}) {
		gameObjectMap := gameObject.(map[string]interface{})
		id := int64(gameObjectMap["id"].(float64))
		gameObjectType := gameObjectMap["type"].(string)
		enabled := gameObjectMap["enabled"].(bool)
		position := physics.Vector2{
			X: gameObjectMap["position"].(map[string]interface{})["x"].(float64),
			Y: gameObjectMap["position"].(map[string]interface{})["y"].(float64),
		}

		if id > uuid {
			uuid = id
		}

		switch gameObjectType {
		case "asteroid":
			asteroid := NewAsteroid(
				id,
				position,
				gameObjectMap["radius"].(float64),
			)
			asteroid.enabled = enabled
			game.manager.AddGameObject(asteroid)
		case "laser":
			fallthrough
		case "rocket":
			owner := game.manager.GetGameObjectByID(int64(gameObjectMap["owner"].(float64)))
			if owner == nil {
				fmt.Println("Owner not found")
				continue
			}
			rotation := gameObjectMap["rotation"].(float64)

			var projectile Projectile
			if gameObjectType == "laser" {
				projectile = *NewLaserProjectile(
					id,
					position,
					rotation,
					owner.(*Spaceship),
				)
			} else {
				projectile = *NewRocketProjectile(
					id,
					position,
					rotation,
					owner.(*Spaceship),
				)
			}

			projectile.enabled = enabled
			projectile.velocity = physics.Vector2{
				X: gameObjectMap["velocity"].(map[string]interface{})["x"].(float64),
				Y: gameObjectMap["velocity"].(map[string]interface{})["y"].(float64),
			}
			projectile.lifespanSec = gameObjectMap["lifespanSec"].(float64)
			projectile.damage = gameObjectMap["damage"].(float64)
			game.manager.AddGameObject(&projectile)
		case "spaceship":
			spaceship := NewSpaceship(
				id,
				gameObjectMap["name"].(string),
				position,
				gameObjectMap["rotation"].(float64),
			)
			spaceship.enabled = enabled
			spaceship.velocity = physics.Vector2{
				X: gameObjectMap["velocity"].(map[string]interface{})["x"].(float64),
				Y: gameObjectMap["velocity"].(map[string]interface{})["y"].(float64),
			}
			spaceship.health = gameObjectMap["health"].(float64)
			spaceship.energy = gameObjectMap["energy"].(float64)
			spaceship.engine.mainThrust = gameObjectMap["engine"].(map[string]interface{})["mainThrust"].(float64)
			spaceship.engine.leftThrust = gameObjectMap["engine"].(map[string]interface{})["leftThrust"].(float64)
			spaceship.engine.rightThrust = gameObjectMap["engine"].(map[string]interface{})["rightThrust"].(float64)
			spaceship.rockets = int32(gameObjectMap["rockets"].(float64))
			spaceship.kills = int32(gameObjectMap["kills"].(float64))
			spaceship.score = gameObjectMap["score"].(float64)
			spaceship.laserReloadTimerSec = gameObjectMap["laserReloadTimerSec"].(float64)
			spaceship.rocketReloadTimerSec = gameObjectMap["rocketReloadTimerSec"].(float64)
			game.manager.AddSpaceship(spaceship)
		case "explosion":
			explosion := NewExplosion(
				id,
				position,
				gameObjectMap["radius"].(float64),
				gameObjectMap["durationSec"].(float64),
			)
			explosion.enabled = enabled
			explosion.lifespanSec = gameObjectMap["lifespanSec"].(float64)
			game.manager.AddGameObject(explosion)
		default:
			fmt.Println("Unknown game object type", gameObjectType)
		}
	}

	// Logs
	logger := game.manager.Logger()
	for _, log := range data["logs"].([]interface{}) {
		logMap := log.(map[string]interface{})
		time, err := time.Parse("2006-01-02 15:04:05", logMap["time"].(string))
		if err != nil {
			fmt.Println(err)
			continue
		}

		id := int64(logMap["id"].(float64))
		if id > uuid {
			uuid = id
		}

		logger.AddMessage(Message{
			id:      id,
			logType: LogType(logMap["logType"].(string)),
			time:    time,
			message: logMap["message"].(string),
			meta:    logMap["meta"].(map[string]interface{}),
		})
	}

	SetUUID(uuid)
	game.status = Status(data["status"].(string))
	return game, nil
}
