package game

import (
	"fmt"
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/davidhorak/space-wars/kernel/utils"
)

func TestSpaceship_SetEngineThrust(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.SetEngineThrust(100, 50, 33)

	if ship.engine.mainThrust != 100 {
		t.Errorf("Expected main thrust %v, got %v", 100, ship.engine.mainThrust)
	}
	if ship.engine.leftThrust != 50 {
		t.Errorf("Expected left thrust %v, got %v", 50, ship.engine.leftThrust)
	}
	if ship.engine.rightThrust != 33 {
		t.Errorf("Expected right thrust %v, got %v", 33, ship.engine.rightThrust)
	}

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
		if err.Error() != test.expectedErr {
			t.Errorf("Expected error %v, got %v", test.expectedErr, err)
		}
	}
}

func TestSpaceship_FireLaser(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.energy = MaxEnergy
	laser, _ := ship.FireLaser(0)

	expectedPosition := physics.Vector2{X: 0, Y: 16}
	if !utils.AlmostEqualVector2(laser.position, expectedPosition) {
		t.Errorf("Expected position %v, got %v", expectedPosition, laser.position)
	}
	if ship.rotation != laser.rotation {
		t.Errorf("Expected rotation %v, got %v", ship.rotation, laser.rotation)
	}
	expectedVelocity := physics.Vector2{X: 0, Y: 320}
	if !utils.AlmostEqualVector2(laser.velocity, expectedVelocity) {
		t.Errorf("Expected velocity %v, got %v", expectedVelocity, laser.velocity)
	}
	if laser.owner != ship {
		t.Errorf("Expected owner %v, got %v", ship, laser.owner)
	}
	expectedColliderPosition := physics.Vector2{X: 0, Y: 16}
	if !utils.AlmostEqualVector2(laser.collider.Position(), expectedColliderPosition) {
		t.Errorf("Expected collider position %v, got %v", expectedColliderPosition, laser.collider.Position())
	}

	ship.energy = 0
	_, err := ship.FireLaser(0)
	if err.Error() != "not enough energy" {
		t.Errorf("Expected error, got %v", err)
	}

	ship.energy = MaxEnergy
	ship.laserReloadTimerSec = 0.1
	_, err = ship.FireLaser(0)
	if err.Error() != "laser is still cooling down" {
		t.Errorf("Expected error, got %v", err)
	}
}

func TestSpaceship_FireRocket(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.energy = MaxEnergy
	ship.rockets = MaxRockets
	rocket, _ := ship.FireRocket(0)

	expectedPosition := physics.Vector2{X: 0, Y: 61}
	if !utils.AlmostEqualVector2(rocket.position, expectedPosition) {
		t.Errorf("Expected position %v, got %v", expectedPosition, rocket.position)
	}
	if ship.rotation != rocket.rotation {
		t.Errorf("Expected rotation %v, got %v", ship.rotation, rocket.rotation)
	}
	expectedVelocity := physics.Vector2{X: 0, Y: 274}
	if !utils.AlmostEqualVector2(rocket.velocity, expectedVelocity) {
		t.Errorf("Expected velocity %v, got %v", expectedVelocity, rocket.velocity)
	}
	if rocket.owner != ship {
		t.Errorf("Expected owner %v, got %v", ship, rocket.owner)
	}
	expectedColliderPosition := physics.Vector2{X: 0, Y: 61}
	if !utils.AlmostEqualVector2(rocket.collider.Position(), expectedColliderPosition) {
		t.Errorf("Expected collider position %v, got %v", expectedColliderPosition, rocket.collider.Position())
	}

	if ship.rockets != MaxRockets-1 {
		t.Errorf("Expected rockets %v, got %v", MaxRockets-1, ship.rockets)
	}
	if ship.rocketReloadTimerSec != RocketReloadSec {
		t.Errorf("Expected rocket shot timer %v, got %v", RocketReloadSec, ship.rocketReloadTimerSec)
	}
	if ship.energy != MaxEnergy-EnergyConsumptionRocket {
		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionRocket, ship.energy)
	}

	ship.energy = 0
	_, err := ship.FireRocket(0)
	if err.Error() != "not enough energy" {
		t.Errorf("Expected error, got %v", err)
	}

	ship.energy = MaxEnergy
	ship.rocketReloadTimerSec = 0.1
	_, err = ship.FireRocket(0)
	if err.Error() != "rocket is not ready to be fired" {
		t.Errorf("Expected error, got %v", err)
	}

	ship.energy = MaxEnergy
	ship.rocketReloadTimerSec = 0
	ship.rockets = 0
	_, err = ship.FireRocket(0)
	if err.Error() != "not enough rockets" {
		t.Errorf("Expected error, got %v", err)
	}
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
		{0, 100, 0, 1, 1, physics.Vector2{X: 31.2, Y: 0}, physics.Vector2{X: 31.2, Y: 0}, 0},
		// 100 right thrust, no main thrust, no drag, 1 tick, 1 second
		{0, 0, 100, 1, 1, physics.Vector2{X: -31.2, Y: 0}, physics.Vector2{X: -31.2, Y: 0}, math.Pi},
		// 100 main thrust, 100 left thrust, no drag, 1 tick, 1 second
		{100, 100, 0, 1, 1, physics.Vector2{X: 31.2, Y: 62.4}, physics.Vector2{X: 31.2, Y: 62.4}, utils.DegreeToRad(63.4349488229)},
		// 100 main thrust, 100 right thrust, no drag, 1 tick, 1 second
		{100, 0, 100, 1, 1, physics.Vector2{X: -31.2, Y: 62.4}, physics.Vector2{X: -31.2, Y: 62.4}, utils.DegreeToRad(116.5650511771)},
	}

	for _, test := range tests {
		ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
		ship.SetEngineThrust(test.mainThrust, test.leftThrust, test.rightThrust)

		for i := 0; i < test.ticks; i++ {
			ship.move(test.deltaTimeSec)
		}

		testCase := fmt.Sprintf("mainThrust: %v, leftThrust: %v, rightThrust: %v, deltaTimeSec: %v, ticks: %v", test.mainThrust, test.leftThrust, test.rightThrust, test.deltaTimeSec, test.ticks)

		if !utils.AlmostEqualVector2(ship.velocity, test.expectedVelocity) {
			t.Log(testCase)
			t.Errorf("Expected velocity %v, got %v", test.expectedVelocity, ship.velocity)
		}
		if !utils.AlmostEqualVector2(ship.position, test.expectedPosition) {
			t.Log(testCase)
			t.Errorf("Expected position %v, got %v", test.expectedPosition, ship.position)
		}
		if !utils.AlmostEqual(ship.rotation, test.expectedRotation) {
			t.Log(testCase)
			t.Errorf("Expected rotation %v, got %v", test.expectedRotation, ship.rotation)
		}
		if !utils.AlmostEqualVector2(ship.collider.Position(), test.expectedPosition) {
			t.Log(testCase)
			t.Errorf("Expected collider position %v, got %v", test.expectedPosition, ship.collider.Position())
		}
	}
}

func TestSpaceship_Move_MaxVelocity(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.SetEngineThrust(100, 0, 0)

	for i := 0; i < 5; i++ {
		ship.move(1)
	}

	if ship.velocity.Magnitude() != MaxVelocitySec {
		t.Errorf("Expected velocity %v, got %v", MaxVelocitySec, ship.velocity.Magnitude())
	}
}

func TestSpaceship_Move_Drag(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.SetEngineThrust(100, 0, 0)

	for i := 0; i < 5; i++ {
		ship.move(1)
	}
	if ship.velocity.Magnitude() != MaxVelocitySec {
		t.Errorf("Expected velocity %v, got %v", MaxVelocitySec, ship.velocity.Magnitude())
	}

	ship.SetEngineThrust(0, 0, 0)

	lastVelocity := ship.velocity.Magnitude()
	for i := 0; i <= 10; i++ {
		ship.move(1)
		newVelocity := ship.velocity.Magnitude()
		if newVelocity > lastVelocity {
			t.Errorf("Expected velocity %v, got %v", lastVelocity, newVelocity)
		}
		lastVelocity = newVelocity
	}

	threshold := float64(MaxVelocitySec / 10)
	if ship.velocity.Magnitude() > threshold {
		t.Errorf("Expected velocity %v, got %v", threshold, ship.velocity.Magnitude())
	}
}

func TestSpaceship_EnergyManagement_Recharge(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.energy = 0

	// +1 to make sure we go over the max
	for i := 0; i < int(MaxEnergy/EnergyRechargeRateSec)+1; i++ {
		ship.energyManagement(1)
	}

	if ship.energy != MaxEnergy {
		t.Errorf("Expected energy %v, got %v", MaxEnergy, ship.energy)
	}
}

func TestSpaceship_EnergyManagement_Thrust(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)

	ship.energy = MaxEnergy
	ship.SetEngineThrust(100, 0, 0)
	ship.energyManagement(1)
	if ship.energy != MaxEnergy-EnergyConsumptionMainThrustSec {
		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionMainThrustSec, ship.energy)
	}

	ship.energy = MaxEnergy
	ship.SetEngineThrust(0, 100, 0)
	ship.energyManagement(1)
	if ship.energy != MaxEnergy-EnergyConsumptionSideThrustSec {
		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionSideThrustSec, ship.energy)
	}

	ship.energy = MaxEnergy
	ship.SetEngineThrust(0, 0, 100)
	ship.energyManagement(1)
	if ship.energy != MaxEnergy-EnergyConsumptionSideThrustSec {
		t.Errorf("Expected energy %v, got %v", MaxEnergy-EnergyConsumptionSideThrustSec, ship.energy)
	}

	ship.energy = MaxEnergy
	ship.SetEngineThrust(50, 50, 50)
	ship.energyManagement(1)
	expected := MaxEnergy - 0.5*(EnergyConsumptionMainThrustSec+EnergyConsumptionSideThrustSec+EnergyConsumptionSideThrustSec)
	if ship.energy != expected {
		t.Errorf("Expected energy %v, got %v", expected, ship.energy)
	}

	ship.energy = MaxEnergy
	ship.SetEngineThrust(100, 0, 0)
	ship.energyManagement(105)
	if ship.energy != 0 {
		t.Errorf("Expected energy %v, got %v", 0, ship.energy)
	}
}

func TestSpaceship_GunManagement(t *testing.T) {
	ship := NewSpaceship(0, "ship", physics.Vector2{X: 0, Y: 0}, math.Pi/2)
	ship.laserReloadTimerSec = LaserReloadSec
	ship.rocketReloadTimerSec = RocketReloadSec

	ship.gunManagement(LaserReloadSec)
	ship.gunManagement(RocketReloadSec)

	if ship.laserReloadTimerSec != 0 {
		t.Errorf("Expected laser shot timer %v, got %v", 0, ship.laserReloadTimerSec)
	}
	if ship.rocketReloadTimerSec != 0 {
		t.Errorf("Expected rocket shot timer %v, got %v", 0, ship.rocketReloadTimerSec)
	}
}
