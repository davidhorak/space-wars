package game

import (
	"fmt"
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/utils"
	"github.com/stretchr/testify/assert"
)

func TestSpaceship_NewSpaceship(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	assert.Equal(t, int64(0), ship.ID())
	assert.Equal(t, true, ship.enabled)
	assert.Equal(t, "ship", ship.name)
	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.position)
	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.startPosition)
	assert.Equal(t, math.Pi/2, ship.rotation)
	assert.Equal(t, math.Pi/2, ship.startRotation)
	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.collider.Position())
	assert.Equal(t, float64(15), ship.collider.Radius())
	assert.Equal(t, physics.Vector2{X: 15, Y: 0}, ship.gunPosition)
	assert.Equal(t, int32(0), ship.kills)
	assert.Equal(t, float64(0), ship.score)
	assert.Equal(t, float64(100), ship.health)
	assert.Equal(t, float64(100), ship.energy)
	assert.Equal(t, int32(10), ship.rockets)
	assert.Equal(t, Engine{
		mainThrust:  0,
		leftThrust:  0,
		rightThrust: 0,
	}, ship.engine)
	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.velocity)
	assert.Equal(t, float64(0), ship.laserReloadTimerSec)
	assert.Equal(t, float64(0), ship.rocketReloadTimerSec)
}

func TestSpaceship_ID(t *testing.T) {
	ship := NewSpaceship(1, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	assert.Equal(t, int64(1), ship.ID())
}

func TestSpaceship_Enabled(t *testing.T) {
	ship := NewSpaceship(1, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	assert.Equal(t, true, ship.Enabled())
}

func TestSpaceship_SetEnabled(t *testing.T) {
	ship := NewSpaceship(1, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.SetEnabled(false)

	assert.Equal(t, false, ship.Enabled())
}

func TestSpaceship_Reset(t *testing.T) {
	ship := NewSpaceship(1, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.Reset()

	assert.Equal(t, true, ship.Enabled())
	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.position)
	assert.Equal(t, math.Pi/2, ship.rotation)
	assert.Equal(t, int32(0), ship.kills)
	assert.Equal(t, float64(0), ship.score)
	assert.Equal(t, float64(100), ship.health)
	assert.Equal(t, float64(100), ship.energy)
	assert.Equal(t, int32(10), ship.rockets)
	assert.Equal(t, Engine{
		mainThrust:  0,
		leftThrust:  0,
		rightThrust: 0,
	}, ship.engine)
	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.velocity)
	assert.Equal(t, float64(0), ship.laserReloadTimerSec)
	assert.Equal(t, float64(0), ship.rocketReloadTimerSec)
}

func TestSpaceship_Position(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	assert.Equal(t, physics.Vector2{X: 0, Y: 0}, ship.Position())
}

func TestSpaceship_SetPosition(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	newPosition := physics.Vector2{X: 15, Y: 25}
	ship.SetPosition(newPosition)

	assert.Equal(t, newPosition, ship.Position())
}

func TestSpaceship_SetStartPosition(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	newPosition := physics.Vector2{X: 15, Y: 25}
	ship.SetStartPosition(newPosition)

	assert.Equal(t, newPosition, ship.startPosition)
}

func TestSpaceship_SetStartRotation(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	newRotation := math.Pi / 3
	ship.SetStartRotation(newRotation)

	assert.Equal(t, newRotation, ship.startRotation)
}

func TestSpaceship_SetEngineThrust(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.SetEngineThrust(100, 50, 33)

	assert.Equal(t, 100.0, ship.engine.mainThrust)
	assert.Equal(t, 50.0, ship.engine.leftThrust)
	assert.Equal(t, 33.0, ship.engine.rightThrust)

	var tests = []struct {
		mainThrust  float64
		leftThrust  float64
		rightThrust float64
		expectedErr string
	}{
		{200, 0, 0, "main thrust must be between 0 and 100"},
		{0, 200, 0, "left thrust must be between 0 and 100"},
		{0, 0, 200, "right thrust must be between 0 and 100"},
		{-100, 0, 0, "main thrust must be between 0 and 100"},
		{0, -100, 0, "left thrust must be between 0 and 100"},
		{0, 0, -100, "right thrust must be between 0 and 100"},
	}

	for _, test := range tests {
		err := ship.SetEngineThrust(test.mainThrust, test.leftThrust, test.rightThrust)
		assert.Contains(t, test.expectedErr, err.Error())
	}
}

func TestSpaceship_FireLaser(t *testing.T) {
	gameManager := NewGameManager()
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.energy = MaxEnergy
	ship.FireLaser(&gameManager)

	assert.Equal(t, 94.0, ship.energy)
	assert.Equal(t, 0.25, ship.laserReloadTimerSec)

	laser := gameManager.gameObjects[0].(*LaserProjectile)

	assert.InDelta(t, 0, laser.position.X, 0.1)
	assert.InDelta(t, 15, laser.position.Y, 0.1)
	assert.InDelta(t, 0, laser.velocity.X, 0.1)
	assert.InDelta(t, 320, laser.velocity.Y, 0.1)
	assert.Equal(t, ship, laser.owner)
	assert.InDelta(t, 0, laser.collider.Position().X, 0.1)
	assert.InDelta(t, 15, laser.collider.Position().Y, 0.1)

	ship.energy = 0
	err := ship.FireLaser(&gameManager)
	assert.Contains(t, "not enough energy", err.Error())

	ship.energy = MaxEnergy
	ship.laserReloadTimerSec = 0.1
	err = ship.FireLaser(&gameManager)
	assert.Contains(t, "laser is still cooling down", err.Error())
}

func TestSpaceship_FireRocket(t *testing.T) {
	gameManager := NewGameManager()
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.energy = MaxEnergy
	ship.FireRocket(&gameManager)

	assert.Equal(t, 80.0, ship.energy)
	assert.Equal(t, 1.0, ship.rocketReloadTimerSec)
	assert.Equal(t, int32(9), ship.rockets)

	rocket := gameManager.gameObjects[0].(*RocketProjectile)

	assert.InDelta(t, 0, rocket.position.X, 0.1)
	assert.InDelta(t, 15, rocket.position.Y, 0.1)
	assert.InDelta(t, 0, rocket.velocity.X, 0.1)
	assert.InDelta(t, 274, rocket.velocity.Y, 0.1)
	assert.Equal(t, ship, rocket.owner)
	assert.InDelta(t, 0, rocket.collider.Position().X, 0.1)
	assert.InDelta(t, 15, rocket.collider.Position().Y, 0.1)

	ship.energy = 0
	err := ship.FireRocket(&gameManager)
	assert.Contains(t, "not enough energy", err.Error())

	ship.energy = MaxEnergy
	ship.laserReloadTimerSec = 0.1
	err = ship.FireRocket(&gameManager)
	assert.Contains(t, "rocket is not ready to be fired", err.Error())

	ship.rockets = 0
	err = ship.FireRocket(&gameManager)
	assert.Contains(t, "not enough rockets", err.Error())
}

func TestSpaceship_Move_Basic(t *testing.T) {
	var tests = []struct {
		mainThrust       float64
		leftThrust       float64
		rightThrust      float64
		deltaTimeSec     float64
		ticks            int
		expectedVelocity physics.Vector2
		expectedPosition physics.Vector2
		expectedRotation float64
	}{
		// 100 main thrust, no side thrust, no drag, 1 tick, 1 second
		{100, 0, 0, 1, 1, physics.Vector2{X: 0, Y: 62.4}, physics.Vector2{X: 0, Y: 62.4}, math.Pi / 2},
		// 100 left thrust, no main thrust, no drag, 1 tick, 1 second
		{0, 100, 0, 1, 1, physics.Vector2{X: -31.2, Y: 0}, physics.Vector2{X: -31.2, Y: 0}, math.Pi},
		// // 100 right thrust, no main thrust, no drag, 1 tick, 1 second
		{0, 0, 100, 1, 1, physics.Vector2{X: 31.2, Y: 0}, physics.Vector2{X: 31.2, Y: 0}, 0},
		// // 100 main thrust, 100 left thrust, no drag, 1 tick, 1 second
		{100, 100, 0, 1, 1, physics.Vector2{X: -31.2, Y: 62.4}, physics.Vector2{X: -31.2, Y: 62.4}, utils.DegreeToRad(116.5650511771)},
		// // 100 main thrust, 100 right thrust, no drag, 1 tick, 1 second
		{100, 0, 100, 1, 1, physics.Vector2{X: 31.2, Y: 62.4}, physics.Vector2{X: 31.2, Y: 62.4}, utils.DegreeToRad(63.4349488229)},
	}

	for _, test := range tests {
		ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
		ship.SetEngineThrust(test.mainThrust, test.leftThrust, test.rightThrust)

		for i := 0; i < test.ticks; i++ {
			ship.move(test.deltaTimeSec)
		}

		testCase := fmt.Sprintf("mainThrust: %v, leftThrust: %v, rightThrust: %v, deltaTimeSec: %v, ticks: %v", test.mainThrust, test.leftThrust, test.rightThrust, test.deltaTimeSec, test.ticks)

		assert.InDelta(t, test.expectedVelocity.X, ship.velocity.X, 0.1, testCase)
		assert.InDelta(t, test.expectedVelocity.Y, ship.velocity.Y, 0.1, testCase)

		assert.InDelta(t, test.expectedPosition.X, ship.position.X, 0.1, testCase)
		assert.InDelta(t, test.expectedPosition.Y, ship.position.Y, 0.1, testCase)
		assert.InDelta(t, test.expectedRotation, ship.rotation, 0.1, testCase)

		assert.InDelta(t, test.expectedPosition.X, ship.collider.Position().X, 0.1, testCase)
		assert.InDelta(t, test.expectedPosition.Y, ship.collider.Position().Y, 0.1, testCase)
	}
}

func TestSpaceship_Move_MaxVelocity(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.SetEngineThrust(100, 0, 0)

	for i := 0; i < 5; i++ {
		ship.move(1)
	}

	assert.InDelta(t, MaxVelocitySec, ship.velocity.Magnitude(), 0.1)
}

func TestSpaceship_Move_Drag(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.SetEngineThrust(100, 0, 0)

	for i := 0; i < 5; i++ {
		ship.move(1)
	}

	assert.InDelta(t, MaxVelocitySec, ship.velocity.Magnitude(), 0.1)

	ship.SetEngineThrust(0, 0, 0)

	for i := 0; i <= 10; i++ {
		ship.move(1)
	}

	threshold := float64(MaxVelocitySec / 10)
	assert.InDelta(t, threshold, ship.velocity.Magnitude(), 0.1)
}

// func TestSpaceship_EnergyManagement_Recharge(t *testing.T) {
// 	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
// 	ship.energy = 0

// 	// +1 to make sure we go over the max
// 	for i := 0; i < int(MaxEnergy/EnergyRechargeRateSec)+1; i++ {
// 		ship.energyManagement(1)
// 	}

// 	if ship.energy != MaxEnergy {
// 		t.Errorf("Expected energy %v, got %v", MaxEnergy, ship.energy)
// 	}
// }

// func TestSpaceship_EnergyManagement_Thrust(t *testing.T) {
// 	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

// 	ship.energy = MaxEnergy
// 	ship.SetEngineThrust(100, 0, 0)
// 	ship.energyManagement(1)
// 	if ship.energy != MaxEnergy-EnergyConsumptionMainThrustSec {
// 		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionMainThrustSec, ship.energy)
// 	}

// 	ship.energy = MaxEnergy
// 	ship.SetEngineThrust(0, 100, 0)
// 	ship.energyManagement(1)
// 	if ship.energy != MaxEnergy-EnergyConsumptionSideThrustSec {
// 		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionSideThrustSec, ship.energy)
// 	}

// 	ship.energy = MaxEnergy
// 	ship.SetEngineThrust(0, 0, 100)
// 	ship.energyManagement(1)
// 	if ship.energy != MaxEnergy-EnergyConsumptionSideThrustSec {
// 		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionSideThrustSec, ship.energy)
// 	}

// 	ship.energy = MaxEnergy
// 	ship.SetEngineThrust(50, 50, 50)
// 	ship.energyManagement(1)
// 	expected := MaxEnergy - 0.5*(EnergyConsumptionMainThrustSec+EnergyConsumptionSideThrustSec+EnergyConsumptionSideThrustSec)
// 	if ship.energy != expected {
// 		t.Errorf("Expected energy %v, got %v", expected, ship.energy)
// 	}

// 	ship.energy = MaxEnergy
// 	ship.SetEngineThrust(100, 0, 0)
// 	ship.energyManagement(105)
// 	if ship.energy != 0 {
// 		t.Errorf("Expected energy %v, got %v", 0, ship.energy)
// 	}
// }

// func TestSpaceship_GunManagement(t *testing.T) {
// 	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
// 	ship.laserReloadTimerSec = LaserReloadSec
// 	ship.rocketReloadTimerSec = RocketReloadSec

// 	ship.gunManagement(LaserReloadSec)
// 	ship.gunManagement(RocketReloadSec)

// 	if ship.laserReloadTimerSec != 0 {
// 		t.Errorf("Expected laser shot timer %v, got %v", 0, ship.laserReloadTimerSec)
// 	}
// 	if ship.rocketReloadTimerSec != 0 {
// 		t.Errorf("Expected rocket shot timer %v, got %v", 0, ship.rocketReloadTimerSec)
// 	}
// }
