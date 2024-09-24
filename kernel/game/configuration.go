package game

const (
	// Asteroid configuration
	MinAsteroids          = 2
	MaxAsteroids          = 7
	MinAsteroidSize       = 10
	MaxAsteroidSize       = 30
	MinAsteroidSeparation = 10 // Minimum distance between asteroids

	// Ship configuration
	ShipSize  = 30
	MaxHealth = 100
	MaxEnergy = 100
	MaxThrust = 100
	// This is tuned to reach edge to edge in 10 seconds for 1920 width
	MaxVelocitySec = 1920 / 10
	// This is tuned to reach the max velocity in 5 seconds
	// with max main thrust, alongside the drag.
	AccelerationCoefficient = 0.325 * MaxVelocitySec
	// This is tuned to reach almost below 10% of max velocity in 10 seconds
	// without any thrust
	DragCoefficient            = 0.2385 * MaxVelocitySec
	SideThrustPowerCoefficient = 0.5 // Relative to the main thrust

	EnergyConsumptionMainThrustSec = MaxEnergy / 8
	EnergyConsumptionSideThrustSec = MaxEnergy / 15
	EnergyRechargeRateSec          = MaxEnergy / 8
	ShipExplosionRadius            = 30
	ShipExplosionDurationSec       = 1
	ScorePerKill                   = 100
	ScorePerDamageCoefficient      = 0.5

	// Laser configuration
	LaserReloadSec         = 0.25
	EnergyConsumptionLaser = 6
	LaserLifespanSec       = 5
	LaserDamage            = 20
	// This is tuned to reach edge to edge in 6 seconds for 1920 width
	LaserVelocitySec          = 1920 / 6
	LaserWidth                = 2
	LaserLength               = 12
	LaserExplosionRadius      = 15
	LaserExplosionDurationSec = 0.75

	// Rocket configuration
	MaxRockets              = 10
	RocketReloadSec         = 1
	EnergyConsumptionRocket = 20
	RocketLifespanSec       = 10
	RocketDamage            = 60
	// This is tuned to reach edge to edge in 7 seconds for 1920 width
	RocketSpeedSec             = 1920 / 7
	RocketDetonateRadius       = 20
	RocketExplosionRadius      = 30
	RocketExplosionDurationSec = 1
)
