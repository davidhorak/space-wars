import { SpaceshipAction, SetEngineThrustAction, FireLaserAction, FireRocketAction } from "../spaceshipAction";
import { SpaceshipManager, SpaceshipManagerFactory, SpaceState } from "../spaceshipManager";
import { Asteroid, Laser, Rocket, Spaceship } from "../types";

const calculateDistance = (x1: number, y1: number, x2: number, y2: number): number => {
    return Math.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2);
}



const willCollide = (x1: number, y1: number, x2: number, y2: number, velocity: {x: number, y: number}, radius: number): boolean => { 
    const objectPos = {xLeft: x1, xRight: x1 + radius, yLeft: y1, yRight: y1 + radius};
    const movePos = {xLeft: x2, xRight: x2 + velocity.x, yLeft: y2, yRight: y2 + velocity.y};

    return !(objectPos.xRight < movePos.xLeft || 
        objectPos.xLeft > movePos.xRight || 
        objectPos.yRight < movePos.yLeft || 
        objectPos.yLeft > movePos.yRight) 
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
    private isRecovering = false;
    // Config
    private ENEGERY_THRESHOLD = 50;
    private DESIRED_ENERGY = 80;
    private ENEGERY_THRESHOLD_LASER = 50;
    private ENEGERY_THRESHOLD_ROCKET = 50;
    private DODGE_THRESHOLD_LASER = 100; // 1 ticks left
    private DODGE_THRESHOLD_ROCKET = 100; // 3 ticks left
    private DODGE_THRESHOLD_ENEMIES = 100; // 1 tick left
    private MOVE_FACTOR = 30;
    private MOVE_ADD = 20;
    // Array of objects
    private asteroidLocs: Array<Asteroid> = [];
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
        // Will enemies hit?
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
        
        const dodgeFactor = this.optimizeDodge() as SetEngineThrustAction[]; 
        const fireLaserAction = this.spaceship.energy > this.ENEGERY_THRESHOLD_LASER && Math.random() < 0.5 ? [["fireLaser"]] : [];
        const fireRocketAction = this.spaceship.energy > this.ENEGERY_THRESHOLD_ROCKET && Math.random() > 0.5 ? [["fireRocket"]]: [];

        console.log(dodgeFactor)

        const dodgeConsumption = this.calculateEnergyConsumption(dodgeFactor[0]?.[1] ?? 0, dodgeFactor[0]?.[2] ?? 0, dodgeFactor[0]?.[3] ?? 0);
        
        return [
            ...dodgeFactor, 
            ...(this.spaceship.energy - dodgeConsumption > this.ENEGERY_THRESHOLD ? fireLaserAction as FireLaserAction[] : []), 
            ...(this.spaceship.energy - dodgeConsumption > this.ENEGERY_THRESHOLD ? fireRocketAction as FireRocketAction[] : [])
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
         // Check if the computed thrust values will cause a collision with any asteroid
        const willCollideWithAsteroidMoveRight = this.asteroidLocs.some(asteroid => {
            const futureX = this.spaceship.position.x + right;
            const futureY = this.spaceship.position.y;

            return willCollide(futureX, futureY, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + + asteroid.collider.radius, { x: 0, y: 0 }, asteroid.radius);
        });

        const willCollideWithAsteroidMoveLeft = this.asteroidLocs.some(asteroid => {
            const futureX = this.spaceship.position.x + left;
            const futureY = this.spaceship.position.y;

            return willCollide(futureX, futureY, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + + asteroid.collider.radius, { x: 0, y: 0 }, asteroid.radius);
        });

        const willCollideWithAsteroidMoveY = this.asteroidLocs.some(asteroid => {
            const futureX = this.spaceship.position.x;
            const futureY = this.spaceship.position.y + main;

            return willCollide(futureX, futureY, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + + asteroid.collider.radius, { x: 0, y: 0 }, asteroid.radius);
        });

        const willCollideWithAsteroidMoveLeftY = this.asteroidLocs.some(asteroid => {
            const futureX = this.spaceship.position.x + left;
            const futureY = this.spaceship.position.y + main;

            return willCollide(futureX, futureY, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + + asteroid.collider.radius, { x: 0, y: 0 }, asteroid.radius);
        })

        const willCollideWithAsteroidMoveRightY = this.asteroidLocs.some(asteroid => {
            const futureX = this.spaceship.position.x + right;
            const futureY = this.spaceship.position.y + main;

            return willCollide(futureX, futureY, asteroid.collider.position.x + asteroid.collider.radius, asteroid.collider.position.y + + asteroid.collider.radius, { x: 0, y: 0 }, asteroid.radius);
        })

        const l = willCollideWithAsteroidMoveLeft || willCollideWithAsteroidMoveLeftY ? right: left > right ? left : 0
        const r = willCollideWithAsteroidMoveRight || willCollideWithAsteroidMoveRightY ? left: right > left ? right : 0
        const lThrust = l != 0 ? l * this.MOVE_FACTOR + this.MOVE_ADD : 0;
        const rThrust = r != 0 ? r * this.MOVE_FACTOR + this.MOVE_ADD : 0; 
        const m = willCollideWithAsteroidMoveY || willCollideWithAsteroidMoveLeftY || willCollideWithAsteroidMoveRightY ? 0: main;
        const mThrust = m != 0 ? m * this.MOVE_FACTOR + this.MOVE_ADD : 0;
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

        
        this.potentialLasersHit.forEach(laser => addAvoidanceVector({x: laser.collider.position.x, y: laser.collider.position.y}));
        this.potentialRocketsHit.forEach(rocket => addAvoidanceVector({x: rocket.collider.position.x, y: rocket.collider.position.y}));
        this.potentialEnemiesHit.forEach(enemy => addAvoidanceVector({x: enemy.collider.position.x, y: enemy.collider.position.y}));
        
        const mainThrust = Math.min(Math.max(avoidanceVector.y, -1), 1);
        const leftThrust = Math.min(Math.max(-avoidanceVector.x, -1), 1);
        const rightThrust = Math.min(Math.max(avoidanceVector.x, -1), 1);

        const [main, left, right] = this.computeMove(mainThrust, leftThrust, rightThrust);
       
      
        const willCollideWithLaser = this.potentialLasersHit.some(laser => {
            const futureX = this.spaceship.position.x + (left > 0 ? left : right);
            const futureY = this.spaceship.position.y + main;
            return willCollide(futureX, futureY, (laser.collider.position.x + laser.collider.size.width) * laser.lifespanSec , (laser.collider.position.y + laser.collider.size.height) * laser.lifespanSec, { x: 0, y: 0 }, 0);
        })
        
        const willCollideWithEnemies = this.potentialEnemiesHit.some(enemy => {
            const futureX = this.spaceship.position.x + (left > 0 ? left : right);
            const futureY = this.spaceship.position.y + main;
            return willCollide(futureX, futureY, enemy.collider.position.x + enemy.collider.radius, enemy.collider.position.y + enemy.collider.radius, { x: 0, y: 0 }, 0);
        })
        
        const willCollideWithRockets = this.potentialRocketsHit.some(rocket => {
            const futureX = this.spaceship.position.x + (left > 0 ? left : right);
            const futureY = this.spaceship.position.y + main;
            return willCollide(futureX, futureY, rocket.collider.position.x + rocket.collider.radius + rocket.velocity.x, rocket.collider.position.y + rocket.collider.radius + rocket.velocity.y, { x: 0, y: 0 }, 0);
        })
        
        console.log(willCollideWithEnemies)
        return willCollideWithLaser || willCollideWithEnemies || willCollideWithRockets || (main == 0 && left == 0 && right == 0) ? [] : [["setEngineThrust", main, left, right]];
    }
}

export const initShip = (): SpaceshipManagerFactory => () =>
    Promise.resolve(new VicecarloansSpaceship());
  