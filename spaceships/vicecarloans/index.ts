import { SpaceshipAction, SetEngineThrustAction, FireLaserAction, FireRocketAction } from "../spaceshipAction";
import { SpaceshipManager, SpaceshipManagerFactory, SpaceState } from "../spaceshipManager";
import { Asteroid, Laser, Rocket, Spaceship } from "../types";

const calculateDistance = (x1: number, y1: number, x2: number, y2: number): number => {
    return Math.abs(Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2));
}

// Predict where the enemy will be based on their current velocity and shoot there
const predictEnemyPosition = (enemy: Spaceship, deltaTime: number): { x: number; y: number } => {
    const futureX = enemy.collider.position.x + (enemy.velocity.x * deltaTime);
    const futureY = enemy.collider.position.y + (enemy.velocity.y * deltaTime);
    return { x: futureX, y: futureY };
};

// Check if there is an asteroid blocking the shot path
const isAsteroidBlockingShot = (ship: Spaceship, enemy: Spaceship, asteroids: Asteroid[]): boolean => {
    const shotLine = { x1: ship.collider.position.x, y1: ship.collider.position.y, x2: enemy.collider.position.x, y2: enemy.collider.position.y };
    return asteroids.some(asteroid => {
        const asteroidX = asteroid.collider.position.x;
        const asteroidY = asteroid.collider.position.y;
        const distanceToShotLine = Math.abs((shotLine.x2 - shotLine.x1) * (shotLine.y1 - asteroidY) - (shotLine.x1 - asteroidX) * (shotLine.y2 - shotLine.y1)) /
            Math.sqrt((shotLine.x2 - shotLine.x1) ** 2 + (shotLine.y2 - shotLine.y1) ** 2);
        return distanceToShotLine < asteroid.collider.radius; // Check if asteroid intersects the shot path
    });
};


class VicecarloansSpaceship implements SpaceshipManager {
    name: string = "vicecarloans"
    private spaceship: Spaceship;
    private width: number;
    private height: number;
    private shipBound = 30
    private attacking = true;
    private isRecovering = false;
    private lastRocketFireMs = 0;
    private lastLaserFireMs = 0;
    private lastDodgeMs = 0;
    // Config
    private ENEGERY_THRESHOLD = 10;
    private DESIRED_ENERGY = 50;
    private DESIRED_ENERGY_FOR_SPIN_MOVE = 30;
    private ENEGERY_THRESHOLD_LASER = 25;
    private ENEGERY_THRESHOLD_ROCKET = 25;
    private DODGE_THRESHOLD_LASER = 200;
    private DODGE_THRESHOLD_ROCKET = 300; 
    private DODGE_THRESHOLD_ENEMIES = 100; 
    private DODGE_THRESHOLD_ASTEROID = 100; 
    private MOVE_FACTOR_Y = 1;
    private MOVE_ADD_Y = 30;
    private MOVE_FACTOR_X = 1;
    private MOVE_ADD_X = 5;
    // Array of objects
    private asteroidLocs: Array<Asteroid> = [];
    private potentialAsteroidsHit: Array<Asteroid> = [];
    private enemies: Array<Spaceship> = [];
    private potentialEnemiesHit: Array<Spaceship> = [];
    private lasers: Array<Laser> = [];
    private potentialLasersHit: Array<Laser> = [];
    private rockets: Array<Rocket> = [];
    private potentialRocketsHit: Array<Rocket> = [];


    onUpdate(state: SpaceState): SpaceshipAction[] {
        if (!this.asteroidLocs.length) {
            this.asteroidLocs = state.gameObjects.filter(go => go.type === "asteroid").map(go => go);
        }
        this.enemies = [];
        this.lasers = [];
        this.rockets = [];
        for(const gameObject of state.gameObjects) {
            if(gameObject.type === "spaceship" && gameObject.name !== this.name) {
                this.enemies.push(gameObject);
            }
            if(gameObject.type === "laser") {
                this.lasers.push(gameObject);
            }
            if(gameObject.type === "rocket") {
                this.rockets.push(gameObject);
            }
        }
        this.spaceship = state.spaceship;
        if(this.spaceship.energy >= this.DESIRED_ENERGY) {
            this.isRecovering = false;
        }
        if(this.spaceship.energy < this.ENEGERY_THRESHOLD) {
            this.isRecovering = true;
        }
        if (this.isRecovering) { 
            return []
        }


        this.potentialAsteroidsHit = this.asteroidLocs.filter(asteroid => {
            return calculateDistance(this.spaceship.collider.position.x, this.spaceship.collider.position.y, asteroid.collider.position.x, asteroid.collider.position.y) < this.DODGE_THRESHOLD_ASTEROID;
        })
        
        this.potentialEnemiesHit = this.enemies.filter(enemy => {
            if (enemy.velocity.x > 0 || enemy.velocity.y > 0) {
                return calculateDistance(this.spaceship.collider.position.x, this.spaceship.collider.position.y, enemy.collider.position.x, enemy.collider.position.y) < this.DODGE_THRESHOLD_ENEMIES;
            }
            return false;
        })
        this.potentialLasersHit = this.lasers.filter(laser => {
            if(this.lasers.length > 0) {
                return calculateDistance(this.spaceship.collider.position.x, this.spaceship.collider.position.y, laser.collider.position.x, laser.collider.position.y) < this.DODGE_THRESHOLD_LASER;
            }
            return false
        })

        this.potentialRocketsHit = this.rockets.filter(rocket => {
            if(this.lasers.length > 0) {
                return calculateDistance(this.spaceship.collider.position.x, this.spaceship.collider.position.y, rocket.collider.position.x, rocket.collider.position.y) < this.DODGE_THRESHOLD_ROCKET;
            }
            return false
        })
        

        // Rocket
        if (this.lastRocketFireMs > 0) {
            this.lastRocketFireMs -= state.deltaTimeMs;
            if (this.lastRocketFireMs <= 0) {
              this.lastRocketFireMs = 0;
            }
        }

        if (this.lastLaserFireMs > 0) {
            this.lastLaserFireMs -= state.deltaTimeMs;
            if (this.lastLaserFireMs <= 0) {
              this.lastLaserFireMs = 0;
            }
        }
        if (this.lastDodgeMs > 0) {
            this.lastDodgeMs -= state.deltaTimeMs;
            if (this.lastDodgeMs <= 0) {
              this.lastDodgeMs = 0;
            }
        }

        if(this.potentialEnemiesHit.length > 1 && this.spaceship.energy > this.DESIRED_ENERGY_FOR_SPIN_MOVE) {
            console.log("PERFORM SPIN MOVE")
            const fireActions: Array<FireLaserAction | FireRocketAction> = []

            if(this.spaceship.rocketReloadTimerSec == 0 && this.lastRocketFireMs <= 0) {
                fireActions.push(["fireRocket"]);
                this.lastRocketFireMs = 500;
            }
            if(this.spaceship.laserReloadTimerSec == 0 && this.lastRocketFireMs <= 0 && this.lastLaserFireMs <= 0) {
                fireActions.push(["fireLaser"]);
            }
            return [
                ["setEngineThrust", 0, 5, 0],
                ...fireActions
            ]
        }

        
        let actions: SpaceshipAction[] = [];
        
        
        if (this.spaceship.energy >= this.DESIRED_ENERGY) {
            actions.push(["setEngineThrust", 0, 5, 0]);

            // Fire rockets if the energy threshold is met and conditions are suitable
            if (this.spaceship.rocketReloadTimerSec === 0 && this.spaceship.rockets > 0 && this.lastRocketFireMs <= 0 && this.lastLaserFireMs <= 0) {
                actions.push(["fireRocket"]);
                this.lastRocketFireMs = 500;
            }
            // Fire lasers based on enemy proximity and aim at predicted enemy position 
            if (this.spaceship.laserReloadTimerSec === 0 && this.lastRocketFireMs <= 0) {
                actions.push(["fireLaser"]);
            }

            
        }

   
        const dodgeActions = this.optimizeDodge()

        if(this.spaceship.energy - this.computeEnergyConsumption(actions) >= this.computeEnergyConsumption(dodgeActions)) {
            actions.push(...dodgeActions);
        }


        if(this.computeEnergyConsumption(actions) > this.spaceship.energy) {
            console.log("NOT ENOUGH ENERGY")
            return []
        } 
        console.log("ENERGY", this.spaceship.energy)
        console.log("ACTIONS", actions)
        return actions
        
    }

    computeEnergyConsumption(actions: SpaceshipAction[]): number {
        let total = 0;
        actions.forEach(action => {
            if (action[0] === "setEngineThrust") {
                total += this.calculateDodgeEnergyConsumption(action[1], action[2], action[3]);
            }
            if (action[0] === "fireLaser") {
                total += 6;
            }
            if (action[0] === "fireRocket") {
                total += 20;
            }
        });
        return total
    }

    calculateDodgeEnergyConsumption(mainThrust: number, leftThrust: number, rightThrust: number): number {
        let energyConsumption = 0
        if(mainThrust > 0) {
            energyConsumption += 12.5
        }
        if(leftThrust > 0 || rightThrust > 0) {
            energyConsumption += 8.33
        }
        
        return energyConsumption
    }

    onStart(spaceship: Spaceship, width: number, height: number): void {
        this.spaceship = spaceship;
        this.width = width;
        this.height = height;
    }

    onReset(spaceship: Spaceship, width: number, height: number): void {
        this.spaceship = spaceship;
        this.width = width;
        this.height = height;
    }

    computeMove(avoidanceVector: {x: number, y: number}): [number, number, number] {
        let rCo = 0
        let lCo = 0
        let mainCo = 0
        if(this.spaceship.rotation > 0) {
            // Headdown
            rCo = avoidanceVector.x < 0 ? Math.abs(avoidanceVector.x * this.MOVE_FACTOR_X) + this.MOVE_ADD_X : 0;
            lCo = avoidanceVector.x > 0 ? Math.abs(avoidanceVector.x * this.MOVE_FACTOR_X) + this.MOVE_ADD_X : 0;
            mainCo = Math.abs(avoidanceVector.y * this.MOVE_FACTOR_Y) + this.MOVE_ADD_Y;
        } else {
            // Head up
            rCo = avoidanceVector.x > 0 ? Math.abs(avoidanceVector.x * this.MOVE_FACTOR_X) + this.MOVE_ADD_X : 0;
            lCo = avoidanceVector.x < 0 ? Math.abs(avoidanceVector.x * this.MOVE_FACTOR_X) + this.MOVE_ADD_X : 0;
            mainCo = Math.abs(avoidanceVector.y * this.MOVE_FACTOR_Y) + this.MOVE_ADD_Y;
        }
        
        const futureXMoveRight = this.spaceship.rotation > 0 ? this.spaceship.position.x + rCo : this.spaceship.position.x - rCo;
        const futureXMoveLeft = this.spaceship.rotation > 0 ? this.spaceship.position.x - lCo : this.spaceship.position.x + lCo;
        const futureYMove = this.spaceship.rotation > 0 ? this.spaceship.position.y + mainCo : this.spaceship.position.y - mainCo;
         // Check if the computed thrust values will cause a collision with any asteroid
        const willCollideWithAsteroidMoveLeft = this.asteroidLocs.some(asteroid => {
            // Move Right?
            return calculateDistance(futureXMoveLeft, this.spaceship.position.y, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + asteroid.collider.radius) < this.DODGE_THRESHOLD_ASTEROID ||
            calculateDistance(futureXMoveLeft, futureYMove, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + asteroid.collider.radius) < this.DODGE_THRESHOLD_ASTEROID; 
        });

        const willCollideWithAsteroidMoveRight = this.asteroidLocs.some(asteroid => {
            return calculateDistance(futureXMoveRight, this.spaceship.position.y, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + asteroid.collider.radius) < this.DODGE_THRESHOLD_ASTEROID ||
            calculateDistance(futureXMoveRight, futureYMove, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + asteroid.collider.radius) < this.DODGE_THRESHOLD_ASTEROID;
        });

        const willCollideWithAsteroidMoveY = this.asteroidLocs.some(asteroid => {
            return calculateDistance(this.spaceship.position.x, futureYMove, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + asteroid.collider.radius) < this.DODGE_THRESHOLD_ASTEROID;
        })

        // Enemies
        const willCollideWithShipMoveLeft = this.potentialEnemiesHit.some(enemy => {
            // Move Right?
            return calculateDistance(futureXMoveLeft, this.spaceship.position.y, enemy.collider.position.x + enemy.collider.radius, enemy.collider.position.y + enemy.collider.radius) < this.DODGE_THRESHOLD_ENEMIES ||
            calculateDistance(futureXMoveLeft, futureYMove, enemy.collider.position.x + enemy.collider.radius, enemy.collider.position.y + enemy.collider.radius) < this.DODGE_THRESHOLD_ENEMIES; 
        });

        const willCollideWithShipMoveRight = this.potentialEnemiesHit.some(enemy => {
            return calculateDistance(futureXMoveRight, this.spaceship.position.y, enemy.collider.position.x + enemy.collider.radius, enemy.collider.position.y + enemy.collider.radius) < this.DODGE_THRESHOLD_ENEMIES ||
            calculateDistance(futureXMoveRight, futureYMove, enemy.collider.position.x + enemy.collider.radius, enemy.collider.position.y + enemy.collider.radius) < this.DODGE_THRESHOLD_ENEMIES;
        });

        const willCollideWithShipMoveY = this.potentialEnemiesHit.some(enemy => {
            return calculateDistance(this.spaceship.position.x, futureYMove, enemy.collider.position.x + enemy.collider.radius, enemy.collider.position.y + enemy.collider.radius) < this.DODGE_THRESHOLD_ENEMIES;
        })


        console.log("LEFT", willCollideWithAsteroidMoveLeft)
        console.log("RIGHT", willCollideWithAsteroidMoveRight)
        console.log("Y", willCollideWithAsteroidMoveY)

        console.log("LEFT SHIP", willCollideWithShipMoveLeft)
        console.log("RIGHT SHIP", willCollideWithShipMoveRight)
        console.log("Y SHIP", willCollideWithShipMoveY)

        let l = lCo
        if (willCollideWithAsteroidMoveLeft) {
            if (this.spaceship.rotation < 0) {
                // Head up
                l = lCo
            } else {
                // Head down
                l = 0
            }
        }

        if(willCollideWithShipMoveLeft) {
            if (this.spaceship.rotation < 0) {
                // Head up
                l = lCo
            } else {
                // Head down
                l = rCo
            }
        }

        
        let r = rCo
        if (willCollideWithAsteroidMoveRight) {
            if (this.spaceship.rotation < 0) {
                // Head up
                r = rCo
            } else {
                // Head down
                r = 0
            }
        }

        if(willCollideWithShipMoveRight) {
            if (this.spaceship.rotation < 0) {
                // Head up
                r = rCo
            } else {
                // Head down
                r = lCo
            }
        }
        
        
        const lThrust = l != 0  ? l : 0;
        const rThrust = r != 0  ? r : 0; 
        const m = willCollideWithAsteroidMoveY ? 0: mainCo;
        const mThrust = m != 0 ? m : 0;
        // return thurst != 0 ? thurst * this.MOVE_FACTOR + 10 : 0;
        return [mThrust, lThrust, rThrust]
    }

    optimizeDodge(): SetEngineThrustAction[] {
        this.lastDodgeMs = 1000;
        const avoidanceVector = { x: 0, y: 0 };
        const addAvoidanceVector = (threat: any, weightFactor: number = 1) => {
            const dx = this.spaceship.position.x - threat.x;
            const dy = this.spaceship.position.y - threat.y;
            const distance = calculateDistance(this.spaceship.position.x, this.spaceship.position.y, threat.x, threat.y);
            if (distance === 0) return; // Avoid division by zero
            const normalizedDx = dx / distance;
            const normalizedDy = dy / distance;
            avoidanceVector.x = Math.max(50, avoidanceVector.x + normalizedDx);
            avoidanceVector.y = Math.max(50, avoidanceVector.y + normalizedDy);
        };

        
        this.potentialLasersHit.forEach(laser => addAvoidanceVector({x: laser.collider.position.x, y: laser.collider.position.y}));
        this.potentialRocketsHit.forEach(rocket => addAvoidanceVector({x: rocket.collider.position.x, y: rocket.collider.position.y}));
        this.potentialEnemiesHit.forEach(enemy => addAvoidanceVector({x: enemy.collider.position.x, y: enemy.collider.position.y}));
        this.asteroidLocs.forEach(asteroid => addAvoidanceVector({x: asteroid.collider.position.x, y: asteroid.collider.position.y}, 2));



        const mainThrust = Math.min(Math.max(avoidanceVector.y, -1), 1);
        const leftThrust = Math.min(Math.max(-avoidanceVector.x, -1), 1);
        const rightThrust = Math.min(Math.max(avoidanceVector.x, -1), 1);

        const [main, left, right] = this.computeMove(avoidanceVector);
       
        
        return main + left + right === 0 ? [] : [["setEngineThrust", main, left, right]];
    }
}

export const initShip = (): SpaceshipManagerFactory => () =>
    Promise.resolve(new VicecarloansSpaceship());
  