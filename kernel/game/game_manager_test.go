package game

import (
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/stretchr/testify/assert"
)

func TestNewGameManager(t *testing.T) {
	manager := NewGameManager()

	assert.Empty(t, manager.gameObjects)
	assert.Empty(t, manager.spaceShips)
	assert.Equal(t, 0, manager.destroyedShips)
	assert.Equal(t, float64(0), manager.gracefulEndTimerMs)
	assert.NotNil(t, manager.logger)
}

func TestGameManager_GameObjects(t *testing.T) {
	manager := NewGameManager()
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)
	spaceship := NewSpaceship(1, "TestShip", physics.Vector2{X: 0, Y: 0}, 100)

	manager.AddGameObject(asteroid)
	manager.AddSpaceship(spaceship)

	assert.ElementsMatch(t, []GameObject{asteroid, spaceship}, manager.GameObjects())
}

func TestGameManager_HasEnded(t *testing.T) {
	manager := NewGameManager()
	ship1 := NewSpaceship(1, "Ship1", physics.Vector2{X: 0, Y: 0}, 100)
	ship2 := NewSpaceship(2, "Ship2", physics.Vector2{X: 0, Y: 0}, 100)

	_ = manager.AddSpaceship(ship1)
	_ = manager.AddSpaceship(ship2)

	assert.False(t, manager.HasEnded(0))

	manager.OnShipDestroyed()

	// Graceful end timer
	assert.False(t, manager.HasEnded(0))
	assert.False(t, manager.HasEnded(ShipExplosionDurationSec*1000+100))
	assert.True(t, manager.HasEnded(0))
}

func TestGameManager_GetGameObjectByIndex(t *testing.T) {
	manager := NewGameManager()
	asteroid := &Asteroid{id: 1}

	manager.AddGameObject(asteroid)
	assert.Equal(t, asteroid, manager.GetGameObjectByIndex(0))
}

func TestGameManager_GameObjectSize(t *testing.T) {
	manager := NewGameManager()
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)
	spaceship := NewSpaceship(1, "TestShip", physics.Vector2{X: 0, Y: 0}, 100)

	assert.Equal(t, 0, manager.GameObjectSize())
	manager.AddGameObject(asteroid)
	manager.AddSpaceship(spaceship)
	assert.Equal(t, 2, manager.GameObjectSize())
}

func TestGameManager_AddGameObject(t *testing.T) {
	manager := NewGameManager()
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)

	manager.AddGameObject(asteroid)
	assert.Equal(t, 1, manager.GameObjectSize())
}

func TestGameManager_AddGameObjects(t *testing.T) {
	manager := NewGameManager()
	asteroid1 := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)
	asteroid2 := NewAsteroid(2, physics.Vector2{X: 0, Y: 0}, 5)

	manager.AddGameObjects([]GameObject{asteroid1, asteroid2})
	assert.Equal(t, 2, manager.GameObjectSize())
}

func TestGameManager_RemoveGameObject(t *testing.T) {
	manager := NewGameManager()
	asteroid := &Asteroid{id: 1}

	manager.AddGameObject(asteroid)
	assert.Equal(t, 1, manager.GameObjectSize())

	manager.RemoveGameObject(asteroid)
	assert.Equal(t, 0, manager.GameObjectSize())
}

func TestGameManager_RemoveGameObjectByIndex(t *testing.T) {
	manager := NewGameManager()
	asteroid1 := &Asteroid{id: 1}
	asteroid2 := &Asteroid{id: 2}

	manager.AddGameObjects([]GameObject{asteroid1, asteroid2})
	assert.Equal(t, 2, manager.GameObjectSize())

	manager.RemoveGameObjectByIndex(0)
	assert.Equal(t, 1, manager.GameObjectSize())
	assert.Equal(t, asteroid2, manager.GetGameObjectByIndex(0))
}

func TestGameManager_AddSpaceship(t *testing.T) {
	manager := NewGameManager()
	ship := NewSpaceship(1, "Ship", physics.Vector2{X: 0, Y: 0}, 100)

	err := manager.AddSpaceship(ship)
	assert.NoError(t, err)
	assert.Len(t, manager.spaceShips, 1)
	assert.Len(t, manager.gameObjects, 1)

	err = manager.AddSpaceship(ship)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "space ship already exists")
}

func TestGameManager_GetSpaceship(t *testing.T) {
	manager := NewGameManager()
	ship := NewSpaceship(1, "Ship", physics.Vector2{X: 0, Y: 0}, 100)

	err := manager.AddSpaceship(ship)
	assert.NoError(t, err)

	ship, err = manager.GetSpaceship("Ship")
	assert.NoError(t, err)
	assert.Equal(t, ship, ship)

	ship, err = manager.GetSpaceship("NonExistentShip")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "space ship not found")
	assert.Nil(t, ship)
}

func TestGameManager_RemoveSpaceship(t *testing.T) {
	manager := NewGameManager()
	ship := NewSpaceship(1, "Ship", physics.Vector2{X: 0, Y: 0}, 100)

	err := manager.AddSpaceship(ship)
	assert.NoError(t, err)

	err = manager.RemoveSpaceship("Ship")
	assert.NoError(t, err)
	assert.Len(t, manager.spaceShips, 0)
	assert.Len(t, manager.gameObjects, 0)

	err = manager.RemoveSpaceship("Ship")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "space ship not found")
}

func TestGameManager_OnShipDestroyed(t *testing.T) {
	manager := NewGameManager()
	ship1 := NewSpaceship(1, "Ship1", physics.Vector2{X: 0, Y: 0}, 100)
	ship2 := NewSpaceship(2, "Ship2", physics.Vector2{X: 0, Y: 0}, 100)

	_ = manager.AddSpaceship(ship1)
	_ = manager.AddSpaceship(ship2)

	manager.OnShipDestroyed()
	assert.Equal(t, 1, manager.destroyedShips)
	assert.Equal(t, float64(ShipExplosionDurationSec*1000+100), manager.gracefulEndTimerMs)

	manager.OnShipDestroyed()
	assert.Equal(t, 2, manager.destroyedShips)
	assert.Equal(t, float64(ShipExplosionDurationSec*1000+100), manager.gracefulEndTimerMs)
}

func TestGameManager_Reset(t *testing.T) {
	manager := NewGameManager()
	ship1 := NewSpaceship(1, "Ship1", physics.Vector2{X: 0, Y: 0}, 100)
	ship2 := NewSpaceship(2, "Ship2", physics.Vector2{X: 0, Y: 0}, 100)
	asteroid := NewAsteroid(1, physics.Vector2{X: 0, Y: 0}, 5)
	explosion := NewExplosion(1, physics.Vector2{X: 0, Y: 0}, 5, 2)

	_ = manager.AddSpaceship(ship1)
	_ = manager.AddSpaceship(ship2)
	manager.AddGameObject(asteroid)
	manager.AddGameObject(explosion)

	manager.Reset()
	assert.Equal(t, 0, manager.destroyedShips)
	assert.Equal(t, float64(0), manager.gracefulEndTimerMs)
	assert.Len(t, manager.spaceShips, 2)
	assert.Len(t, manager.gameObjects, 3)
}
