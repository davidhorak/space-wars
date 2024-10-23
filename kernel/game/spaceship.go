package game

import (
	"errors"
	"math"
	"time"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/physics/collider"
)

type Engine struct {
	mainThrust  float64 // 0-100
	leftThrust  float64 // 0-100
	rightThrust float64 // 0-100
}

type Spaceship struct {
	id                   int64
	name                 string
	enabled              bool
	rotation             float64
	startRotation        float64
	position             physics.Vector2
	startPosition        physics.Vector2
	gunPosition          physics.Vector2 // Relative to the ship's position, orientation to rad 0
	velocity             physics.Vector2
	health               float64 // 0-100
	energy               float64 // 0-100
	engine               Engine
	rockets              int32
	kills                int32
	score                float64
	collider             collider.CircleCollider
	laserReloadTimerSec  float64
	rocketReloadTimerSec float64
}

func NewSpaceship(id int64, name string, position physics.Vector2, rotation float64) *Spaceship {
	ship := &Spaceship{
		id:            id,
		enabled:       true,
		name:          name,
		position:      position,
		startPosition: position,
		rotation:      rotation,
		startRotation: rotation,
		// TODO: Create polygon collider
		collider:    *collider.NewCircleCollider(position, ShipSize/2),
		gunPosition: physics.Vector2{X: ShipSize / 2, Y: 0},
	}
	ship.Reset()
	return ship
}

func (ship *Spaceship) ID() int64 {
	return ship.id
}

func (ship *Spaceship) Enabled() bool {
	return ship.enabled
}

func (ship *Spaceship) SetEnabled(enabled bool) {
	ship.enabled = enabled
}

func (ship *Spaceship) Reset() {
	ship.enabled = true
	ship.position = ship.startPosition
	ship.rotation = ship.startRotation
	ship.health = MaxHealth
	ship.energy = MaxEnergy
	ship.rockets = MaxRockets
	ship.engine = Engine{
		mainThrust:  0,
		leftThrust:  0,
		rightThrust: 0,
	}
	ship.kills = 0
	ship.score = 0
	ship.velocity = physics.Vector2{
		X: 0,
		Y: 0,
	}
	ship.laserReloadTimerSec = 0
	ship.rocketReloadTimerSec = 0
}

func (ship *Spaceship) Position() physics.Vector2 {
	return ship.position
}

func (ship *Spaceship) SetPosition(position physics.Vector2) {
	ship.position = position
}

func (ship *Spaceship) SetStartPosition(position physics.Vector2) {
	ship.startPosition = position
}

func (ship *Spaceship) SetStartRotation(rotation float64) {
	ship.startRotation = rotation
}

func (ship *Spaceship) Update(deltaTimeMs float64, gameManager *GameManager) {
	deltaTimeSec := deltaTimeMs / 1000

	ship.gunManagement(deltaTimeSec)
	ship.energyManagement(deltaTimeSec)
	if ship.energy <= 0 {
		ship.SetEngineThrust(0, 0, 0)
		return
	}
	ship.move(deltaTimeSec)
}

func (ship *Spaceship) Collider() collider.Collider {
	return &ship.collider
}

func (ship *Spaceship) SetEngineThrust(main, left, right float64) error {
	if main < 0 || main > 100 {
		return errors.New("main thrust must be between 0 and 100")
	}
	if left < 0 || left > 100 {
		return errors.New("left thrust must be between 0 and 100")
	}
	if right < 0 || right > 100 {
		return errors.New("right thrust must be between 0 and 100")
	}

	ship.engine.mainThrust = main
	ship.engine.leftThrust = left
	ship.engine.rightThrust = right
	return nil
}

func (ship *Spaceship) FireLaser(gameManager *GameManager) error {
	if ship.energy < EnergyConsumptionLaser {
		return errors.New("not enough energy")
	}

	if ship.laserReloadTimerSec > 0 {
		return errors.New("laser is still cooling down")
	}

	ship.energy -= EnergyConsumptionLaser
	ship.laserReloadTimerSec = LaserReloadSec

	gameManager.AddGameObject(NewLaserProjectile(
		NewUUID(),
		ship.position.Add(ship.gunPosition.Rotate(ship.rotation)),
		ship.rotation,
		ship,
	))
	return nil
}

func (ship *Spaceship) FireRocket(gameManager *GameManager) error {
	if ship.rockets == 0 {
		return errors.New("not enough rockets")
	}
	if ship.energy < EnergyConsumptionRocket {
		return errors.New("not enough energy")
	}
	if ship.rocketReloadTimerSec > 0 {
		return errors.New("rocket is not ready to be fired")
	}

	ship.rockets--
	ship.energy -= EnergyConsumptionRocket
	ship.rocketReloadTimerSec = RocketReloadSec
	gameManager.AddGameObject(NewRocketProjectile(
		NewUUID(),
		ship.position.Add(ship.gunPosition.Rotate(ship.rotation)),
		ship.rotation,
		ship,
	))
	return nil
}

func (ship *Spaceship) HasKilled(target *Spaceship) {
	ship.kills++
	ship.score += ScorePerKill
}

func (ship *Spaceship) OnCollision(other GameObject, gameManager *GameManager, order int) {
	switch other.(type) {
	case *Asteroid:
		ship.TakeDamage(MaxHealth, gameManager, nil)
		gameManager.Logger().Collision(time.Now(), ship.name, "an asteroid")
	case *Spaceship:
		ship.TakeDamage(MaxHealth, gameManager, nil)
		if order == 0 {
			gameManager.Logger().Collision(time.Now(), ship.name, other.(*Spaceship).name)
		}
	default:
		return
	}
}

func (ship *Spaceship) TakeDamage(damage float64, gameManager *GameManager, damageDealer *Spaceship) {
	ship.health -= damage
	ship.health = math.Max(ship.health, 0)
	if ship.health <= 0 {
		ship.destroy(gameManager)
		if damageDealer != nil {
			gameManager.Logger().Kill(time.Now(), ship.name, damageDealer.name)
			damageDealer.HasKilled(ship)
		}
	}
}

func (ship *Spaceship) AddScore(score float64) {
	ship.score += score
}

func (ship *Spaceship) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":      "spaceship",
		"id":        ship.id,
		"enabled":   ship.enabled,
		"destroyed": ship.health <= 0,
		"name":      ship.name,
		"startPosition": map[string]interface{}{
			"x": ship.startPosition.X,
			"y": ship.startPosition.Y,
		},
		"position": map[string]interface{}{
			"x": ship.position.X,
			"y": ship.position.Y,
		},
		"rotation": ship.rotation,
		"velocity": map[string]interface{}{
			"x": ship.velocity.X,
			"y": ship.velocity.Y,
		},
		"health": ship.health,
		"energy": ship.energy,
		"engine": map[string]interface{}{
			"mainThrust":  ship.engine.mainThrust,
			"leftThrust":  ship.engine.leftThrust,
			"rightThrust": ship.engine.rightThrust,
		},
		"rockets":              ship.rockets,
		"kills":                ship.kills,
		"score":                ship.score,
		"laserReloadTimerSec":  ship.laserReloadTimerSec,
		"rocketReloadTimerSec": ship.rocketReloadTimerSec,
		"collider":             ship.collider.Serialize(),
		// TODO: Add collider, if polygon
	}
}

func (ship *Spaceship) gunManagement(deltaTimeSec float64) {
	ship.laserReloadTimerSec -= deltaTimeSec
	ship.rocketReloadTimerSec -= deltaTimeSec

	if ship.laserReloadTimerSec < 0 {
		ship.laserReloadTimerSec = 0
	}
	if ship.rocketReloadTimerSec < 0 {
		ship.rocketReloadTimerSec = 0
	}
}

func (ship *Spaceship) energyManagement(deltaTimeSec float64) {
	// TODO: Investigate if this is needed
	// if ship.engine.mainThrust == 0 && ship.engine.leftThrust == 0 && ship.engine.rightThrust == 0 {
	ship.energy += deltaTimeSec * EnergyRechargeRateSec
	ship.energy = math.Min(ship.energy, MaxEnergy)
	// 	return
	// }

	ship.energy -= ship.engine.mainThrust / MaxThrust * deltaTimeSec * EnergyConsumptionMainThrustSec
	ship.energy -= ship.engine.leftThrust / MaxThrust * deltaTimeSec * EnergyConsumptionSideThrustSec
	ship.energy -= ship.engine.rightThrust / MaxThrust * deltaTimeSec * EnergyConsumptionSideThrustSec
	ship.energy = math.Max(ship.energy, 0)
}

func (ship *Spaceship) move(deltaTimeSec float64) {
	direction := physics.Vector2{X: 1, Y: 0}

	mainThrust := direction.Rotate(ship.rotation)
	mainThrust = mainThrust.Multiply(ship.engine.mainThrust / MaxThrust * deltaTimeSec * AccelerationCoefficient)
	leftThrust := direction.Rotate(ship.rotation + math.Pi/2)
	leftThrust = leftThrust.Multiply(ship.engine.leftThrust / MaxThrust * SideThrustPowerCoefficient * deltaTimeSec * AccelerationCoefficient)
	rightThrust := direction.Rotate(ship.rotation - math.Pi/2)
	rightThrust = rightThrust.Multiply(ship.engine.rightThrust / MaxThrust * SideThrustPowerCoefficient * deltaTimeSec * AccelerationCoefficient)

	drag := direction.Rotate(ship.rotation + math.Pi)
	// TODO: investigate if this should be divided by deltaTimeSec
	drag = drag.Multiply(ship.velocity.Magnitude() / MaxVelocitySec * deltaTimeSec * DragCoefficient)

	ship.velocity = ship.velocity.Add(mainThrust)
	ship.velocity = ship.velocity.Add(leftThrust)
	ship.velocity = ship.velocity.Add(rightThrust)
	ship.velocity = ship.velocity.Add(drag)
	ship.velocity = ship.velocity.Clamp(MaxVelocitySec / deltaTimeSec)

	ship.position = ship.position.Add(ship.velocity.Multiply(deltaTimeSec))
	if ship.velocity.Magnitude() > 0 {
		ship.rotation = math.Atan2(ship.velocity.Y, ship.velocity.X)
	}

	ship.collider.SetPosition(ship.position)
	// TODO: Apply rotation for the polygon collider
}

func (ship *Spaceship) destroy(gameManager *GameManager) {
	ship.enabled = false
	gameManager.AddGameObject(NewExplosion(
		NewUUID(),
		physics.Vector2{
			X: ship.position.X - float64(ShipExplosionRadius),
			Y: ship.position.Y - float64(ShipExplosionRadius),
		},
		float64(ShipExplosionRadius),
		float64(ShipExplosionDurationSec),
	))
	gameManager.OnShipDestroyed()
}
