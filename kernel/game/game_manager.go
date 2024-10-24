package game

import "fmt"

type GameManager struct {
	gameObjects        []GameObject
	spaceShips         map[string]*Spaceship
	destroyedShips     int
	gracefulEndTimerMs float64
	logger             Logger
}

func NewGameManager() GameManager {
	return GameManager{
		gameObjects:    []GameObject{},
		spaceShips:     map[string]*Spaceship{},
		logger:         NewLogger(),
		destroyedShips: 0,
	}
}

func (manager *GameManager) GameObjects() []GameObject {
	return manager.gameObjects
}

func (manager *GameManager) HasEnded(deltaTimeMs float64) bool {
	if manager.gracefulEndTimerMs > 0 {
		manager.gracefulEndTimerMs -= deltaTimeMs
		return false
	}
	return manager.destroyedShips >= len(manager.spaceShips)-1
}

func (manager *GameManager) GetGameObjectByID(id int64) GameObject {
	for _, gameObject := range manager.gameObjects {
		if gameObject.ID() == id {
			return gameObject
		}
	}
	return nil
}

func (manager *GameManager) GetGameObjectByIndex(index int) GameObject {
	return manager.gameObjects[index]
}

func (manager *GameManager) GameObjectSize() int {
	return len(manager.gameObjects)
}

func (manager *GameManager) AddGameObject(gameObject GameObject) {
	manager.gameObjects = append(manager.gameObjects, gameObject)
}

func (manager *GameManager) AddGameObjects(gameObjects []GameObject) {
	manager.gameObjects = append(manager.gameObjects, gameObjects...)
}

func (manager *GameManager) RemoveGameObject(gameObject GameObject) {
	for i, obj := range manager.gameObjects {
		if obj.ID() == gameObject.ID() {
			manager.gameObjects = append(manager.gameObjects[:i], manager.gameObjects[i+1:]...)
			break
		}
	}
}

func (manager *GameManager) RemoveGameObjectByIndex(index int) {
	manager.gameObjects = append(manager.gameObjects[:index], manager.gameObjects[index+1:]...)
}

func (manager *GameManager) AddSpaceship(spaceShip *Spaceship) error {
	if _, ok := manager.spaceShips[spaceShip.name]; ok {
		return fmt.Errorf("space ship already exists: %s", spaceShip.name)
	}

	manager.spaceShips[spaceShip.name] = spaceShip
	manager.gameObjects = append(manager.gameObjects, spaceShip)
	return nil
}

func (manager *GameManager) GetSpaceship(name string) (*Spaceship, error) {
	spaceShip, ok := manager.spaceShips[name]
	if !ok {
		return nil, fmt.Errorf("space ship not found: %s", name)
	}
	return spaceShip, nil
}

func (manager *GameManager) RemoveSpaceship(name string) error {
	spaceShip, err := manager.GetSpaceship(name)
	if err != nil {
		return err
	}

	manager.RemoveGameObject(spaceShip)
	delete(manager.spaceShips, name)
	return nil
}

func (manager *GameManager) OnShipDestroyed() {
	manager.destroyedShips++

	shipsToScore := 1
	if len(manager.spaceShips) > 3 {
		shipsToScore = 3
	}
	if manager.destroyedShips >= len(manager.spaceShips)-shipsToScore {
		for _, spaceShip := range manager.spaceShips {
			if spaceShip.enabled {
				spaceShip.AddScore(ScorePerSurvivor)
			}
		}
	}

	if manager.destroyedShips >= len(manager.spaceShips)-1 {
		manager.gracefulEndTimerMs = (ShipExplosionDurationSec * 1000) + 100
}

func (manager *GameManager) Reset() {
	gameObjects := make([]GameObject, 0)
	for _, gameObject := range manager.GameObjects() {
		switch gameObject.(type) {
		case *Spaceship:
			gameObject.(*Spaceship).Reset()
			gameObjects = append(gameObjects, gameObject)
		case *Asteroid:
			gameObjects = append(gameObjects, gameObject)
		}
	}
	manager.gameObjects = gameObjects
	manager.destroyedShips = 0
	manager.gracefulEndTimerMs = 0
}

func (manager *GameManager) Logger() Logger {
	return manager.logger
}
