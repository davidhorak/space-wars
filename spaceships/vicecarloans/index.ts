import { SpaceshipAction, SetEngineThrustAction, FireLaserAction, FireRocketAction } from "../spaceshipAction";
import { SpaceshipManager, SpaceshipManagerFactory, SpaceState } from "../spaceshipManager";
import { Asteroid, Laser, Rocket, Spaceship } from "../types";

const calculateDistance = (x1: number, y1: number, x2: number, y2: number): number => {
    return Math.abs(Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2));
}

const willCollide = (x1: number, y1: number, x2: number, y2: number, velocity: {x: number, y: number}, radius: number): boolean => { 
    const objectPos = {xLeft: x2, xRight: x2 + radius, yLeft: y2, yRight: y2 + radius};
    const movePos = {xLeft: x1, xRight: x1 + velocity.x, yLeft: y1, yRight: y1 + velocity.y};

    return objectPos.xRight < movePos.xLeft &&
        objectPos.xLeft > movePos.xRight && 
        objectPos.yRight < movePos.yLeft && 
        objectPos.yLeft > movePos.yRight
}
 
const findMinimumCollisionVelocity = (x1: number, y1: number, x2: number, y2: number, radius: number): {x: number, y: number} | null => {
    const velocityRange = 100; // Define the range of velocities to test
    const step = 1; // Define the step size for velocity increments

    for (let vx = -velocityRange; vx <= velocityRange; vx += step) {
        for (let vy = -velocityRange; vy <= velocityRange; vy += step) {
            if (willCollide(x1, y1, x2, y2, {x: vx, y: vy}, radius)) {
                return {x: vx, y: vy}; // Return the first velocity that causes a collision
            }
        }
    }

    return null; // Return null if no collision velocity is found
}

class VicecarloansSpaceship implements SpaceshipManager {
    name: string = "vicecarloans"
    private spaceship: Spaceship;
    private width: number;
    private height: number;
    private shipBound = 30
    private attacking = true;
    private isRecovering = false;
    private lastRocketFireMs = 0;
    // Config
    private ENEGERY_THRESHOLD = 10;
    private DESIRED_ENERGY = 30;
    private ENEGERY_THRESHOLD_LASER = 25;
    private ENEGERY_THRESHOLD_ROCKET = 25;
    private DODGE_THRESHOLD_LASER = 200;
    private DODGE_THRESHOLD_ROCKET = 300; 
    private DODGE_THRESHOLD_ENEMIES = 100; 
    private DODGE_THRESHOLD_ASTEROID = 50; 
    private MOVE_FACTOR_Y = 30;
    private MOVE_ADD_Y = 40;
    private MOVE_FACTOR_X = 30;
    private MOVE_ADD_X = 20;
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

        
        const fireRocketAction = this.spaceship.energy > this.ENEGERY_THRESHOLD_ROCKET && this.spaceship.rocketReloadTimerSec == 0 && this.spaceship.rockets > 0 && Math.random() > 0.5 ? [["fireRocket"]]: [];

        if(fireRocketAction.length > 0) {
            this.lastRocketFireMs = 500;
        }
        const fireLaserAction = this.spaceship.energy > 10 && this.lastRocketFireMs <= 0 && this.spaceship.energy > this.ENEGERY_THRESHOLD_LASER && this.spaceship.laserReloadTimerSec == 0 && Math.random() < 0.5 ? [["fireLaser"]] : [];

        const dodgeFactor = (this.potentialEnemiesHit.length || this.potentialLasersHit.length || this.potentialRocketsHit.length || this.potentialAsteroidsHit) ? this.optimizeDodge() as SetEngineThrustAction[] : []

        
        console.log(dodgeFactor)

        
        return [
            ...dodgeFactor, 
            ...(fireRocketAction as FireRocketAction[]),
            ...(fireLaserAction as FireLaserAction[]), 
           
        ];
    }

    calculateEnergyConsumption(mainThrust: number, leftThrust: number, rightThrust: number): number {
        // Example energy consumption calculation
        return Math.abs(mainThrust) + Math.abs(leftThrust) + Math.abs(rightThrust);
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

    computeMove(main: number, left: number, right: number): [number, number, number] {
        console.log(main, left, right)
        const futureXMoveRight = this.spaceship.position.x + right * this.MOVE_FACTOR_X + this.MOVE_ADD_X;
        const futureXMoveLeft = this.spaceship.position.x + left * this.MOVE_FACTOR_X + this.MOVE_ADD_X;
        const futureYMove = this.spaceship.position.y + main * this.MOVE_FACTOR_Y + this.MOVE_ADD_Y;
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

        const l = willCollideWithAsteroidMoveLeft || willCollideWithShipMoveLeft ? 0: left
        const r = willCollideWithAsteroidMoveRight || willCollideWithShipMoveRight? 0: right
        const lThrust = l != 0 ? l * this.MOVE_FACTOR_X + this.MOVE_ADD_X : 0;
        const rThrust = r != 0 ? r * this.MOVE_FACTOR_X + this.MOVE_ADD_X : 0; 
        const m = willCollideWithAsteroidMoveY || willCollideWithShipMoveY ? 0: main;
        const mThrust = m != 0 ? m * this.MOVE_FACTOR_Y + this.MOVE_ADD_Y : 0;
        // return thurst != 0 ? thurst * this.MOVE_FACTOR + 10 : 0;
        return [mThrust, lThrust, rThrust]
    }

    optimizeDodge(): SpaceshipAction[] {
        const avoidanceVector = { x: 0, y: 0 };
        const addAvoidanceVector = (threat: any, weightFactor: number = 1) => {
            const dx = this.spaceship.position.x - threat.x;
            const dy = this.spaceship.position.y - threat.y;
            const distance = calculateDistance(this.spaceship.position.x, this.spaceship.position.y, threat.x, threat.y);
            if (distance === 0) return; // Avoid division by zero
            const normalizedDx = dx / distance;
            const normalizedDy = dy / distance;
            avoidanceVector.x += normalizedDx;
            avoidanceVector.y += normalizedDy;
        };

        console.log(this.potentialAsteroidsHit)

        
        this.potentialLasersHit.forEach(laser => addAvoidanceVector({x: laser.collider.position.x, y: laser.collider.position.y}));
        this.potentialRocketsHit.forEach(rocket => addAvoidanceVector({x: rocket.collider.position.x, y: rocket.collider.position.y}));
        this.potentialEnemiesHit.forEach(enemy => addAvoidanceVector({x: enemy.collider.position.x, y: enemy.collider.position.y}));
        this.potentialAsteroidsHit.forEach(asteroid => addAvoidanceVector({x: asteroid.collider.position.x, y: asteroid.collider.position.y}));
        
        const mainThrust = Math.min(Math.max(avoidanceVector.y, -1), 1);
        const leftThrust = Math.min(Math.max(-avoidanceVector.x, -1), 1);
        const rightThrust = Math.min(Math.max(avoidanceVector.x, -1), 1);

        const [main, left, right] = this.computeMove(mainThrust, leftThrust, rightThrust);
       
      
        
        return (main == 0 && left == 0 && right == 0) ? [] : [["setEngineThrust", main, left, right]];
    }
}

export const initShip = (): SpaceshipManagerFactory => () =>
    Promise.resolve(new VicecarloansSpaceship());
  