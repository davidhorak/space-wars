import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
  Spaceship,
  GameObject,
} from "..";
import { SetEngineThrustAction } from "../spaceshipAction";
import { Projectile } from "../types";

const LASER_SPEED = 320;
const ROCKET_SPEED = 274;

const GAME_WIDTH = 1024;
const GAME_HEIGHT = 768;

const LASER_LIFESPAN = 5;
const ROCKET_LIFESPAN = 10;

class KodySpaceshipManager implements SpaceshipManager {
  name = "kody";

  startTs: number;

  elapsedTime() {
    return Date.now() - this.startTs;
  }

  /**
   *
   * @param ship MY ship
   * @param target Target ship
   * @param speed Speed of the projectile
   * @param lifespan Lifespan of the projectile
   * @returns
   */
  // Margin for projectile hit detection
  private PROJECTILE_HIT_MARGIN = 5; // Increased from 3 to 5 units

  shouldFire(
    ship: Spaceship,
    target: Spaceship,
    speed: number,
    lifespan: number
  ) {
    // Calculate the projectile's velocity components based on ship's rotation
    const projectileVx = speed * Math.cos(ship.rotation);
    const projectileVy = speed * Math.sin(ship.rotation);

    // Calculate the relative position and velocity
    const relativeX = target.position.x - ship.position.x;
    const relativeY = target.position.y - ship.position.y;
    const relativeVx = target.velocity.x - projectileVx;
    const relativeVy = target.velocity.y - projectileVy;

    // Calculate coefficients for the quadratic equation
    const a = relativeVx * relativeVx + relativeVy * relativeVy;
    const b = 2 * (relativeX * relativeVx + relativeY * relativeVy);
    const c =
      relativeX * relativeX +
      relativeY * relativeY -
      Math.pow(this.PROJECTILE_HIT_MARGIN, 2);

    // Solve the quadratic equation
    const discriminant = b * b - 4 * a * c;
    if (discriminant < 0) {
      return false; // No real solutions, projectile won't come within the margin
    }

    const t1 = (-b + Math.sqrt(discriminant)) / (2 * a);
    const t2 = (-b - Math.sqrt(discriminant)) / (2 * a);

    // Check if either solution is within the projectile's lifespan
    if ((t1 > 0 && t1 <= lifespan) || (t2 > 0 && t2 <= lifespan)) {
      return true;
    }

    return false;
  }

  willCollide(ship: Spaceship, obj: Projectile | Spaceship) {
    /**
     * ship is me
     * obj is the other object
     *
     * both myself and the obj has position and velocity
     *
     * calculate if i will intersect with the object within one second
     *
     * in this calculation have a togglable collision radius
     *
     * return true if i will intersect with the object within one second
     */
    const COLLISION_RADIUS = 25; // Adjustable collision radius
    const TIME_HORIZON = 1; // Check for collisions within 1 second

    // Calculate relative position and velocity
    const relativeX = obj.position.x - ship.position.x;
    const relativeY = obj.position.y - ship.position.y;
    const relativeVx = obj.velocity.x - ship.velocity.x;
    const relativeVy = obj.velocity.y - ship.velocity.y;

    // Calculate coefficients for the quadratic equation
    const a = relativeVx * relativeVx + relativeVy * relativeVy;
    const b = 2 * (relativeX * relativeVx + relativeY * relativeVy);
    const c =
      relativeX * relativeX +
      relativeY * relativeY -
      Math.pow(COLLISION_RADIUS, 2);

    // Solve the quadratic equation
    const discriminant = b * b - 4 * a * c;
    if (discriminant < 0) {
      return false; // No real solutions, objects won't collide
    }

    const t1 = (-b + Math.sqrt(discriminant)) / (2 * a);
    const t2 = (-b - Math.sqrt(discriminant)) / (2 * a);

    // Check if either solution is within the time horizon
    if ((t1 > 0 && t1 <= TIME_HORIZON) || (t2 > 0 && t2 <= TIME_HORIZON)) {
      return true;
    }

    return false;
  }

  calcMovement(state: SpaceState): SetEngineThrustAction {
    /**
     * state.spaceship is me
     * state.gameObjects is all other objects
     *
     * The objective is to dodge everything.
     *
     * First, iterate through every game object (skipping myself)
     *  if the object is a laser or a rocket, then check if it will run into me
     *    if it will, turn and go forward
     * if the object is a spaceship, turn towards it
     * otherwise, turn in place
     */

    const ship = state.spaceship;
    let mainThrust = 0;
    let leftThrust = 0;
    let rightThrust = 0;
    let targetRotation = ship.rotation;

    let willCollide = false;
    for (const obj of state.gameObjects) {
      if (obj.id === ship.id) continue;

      if (obj.type === "laser" || obj.type === "rocket") {
        if (this.willCollide(ship, obj)) {
          willCollide = true;
          // Dodge by turning perpendicular to the incoming projectile
          const angleToProjectile = Math.atan2(
            obj.position.y - ship.position.y,
            obj.position.x - ship.position.x
          );
          targetRotation = angleToProjectile + Math.PI / 2;
          mainThrust = 100; // Full thrust to dodge
          break;
        }
      }
    }

    if (willCollide) {
      // Adjust rotation
      const rotationDiff =
        ((targetRotation - ship.rotation + Math.PI * 3) % (Math.PI * 2)) -
        Math.PI;
      if (Math.abs(rotationDiff) > 0.1) {
        if (rotationDiff > 0) {
          rightThrust = 100;
        } else {
          leftThrust = 100;
        }
      }

      return ["setEngineThrust", mainThrust, leftThrust, rightThrust];
    }

    // Find the closest spaceship
    let closestShip: Spaceship | null = null;
    let closestDistance = Infinity;

    for (const obj of state.gameObjects) {
      if (obj.type === "spaceship" && obj.id !== ship.id && !obj.destroyed) {
        const distance = Math.hypot(
          obj.position.x - ship.position.x,
          obj.position.y - ship.position.y
        );
        if (distance < closestDistance) {
          closestDistance = distance;
          closestShip = obj as Spaceship;
        }
      }
    }

    if (closestShip) {
      // Calculate future position of the closest ship
      const predictionTime = 1; // Adjustable prediction time in seconds
      const futureX =
        closestShip.position.x + closestShip.velocity.x * predictionTime;
      const futureY =
        closestShip.position.y + closestShip.velocity.y * predictionTime;

      // Calculate angle to the predicted future position of the closest ship
      const angleToShip = Math.atan2(
        futureY - ship.position.y,
        futureX - ship.position.x
      );

      // Calculate the difference between current rotation and desired rotation
      const rotationDiff =
        ((angleToShip - ship.rotation + Math.PI * 3) % (Math.PI * 2)) - Math.PI;

      // Determine which direction to rotate
      if (Math.abs(rotationDiff) > 0.1) {
        if (rotationDiff > 0) {
          leftThrust = 100;
        } else {
          rightThrust = 100;
        }
      }
    }
    return ["setEngineThrust", 0, leftThrust, rightThrust];
  }

  onUpdate(state: SpaceState): SpaceshipAction[] {
    // lazer speed 320
    // rocket speed 274
    let shouldFireLaser = false;
    let shouldFireRocket = false;

    for (const obj of state.gameObjects) {
      if (
        obj.type !== "spaceship" ||
        obj.id === state.spaceship.id ||
        obj.destroyed
      ) {
        continue;
      }

      // if (
      //   this.shouldFire(state.spaceship, obj, LASER_SPEED, LASER_LIFESPAN)
      // ) {
      //   shouldFireLaser = true;
      // }

      if (
        this.shouldFire(state.spaceship, obj, ROCKET_SPEED, ROCKET_LIFESPAN)
      ) {
        shouldFireRocket = true;
      }
    }

    const resp: SpaceshipAction[] = [this.calcMovement(state)];

    if (shouldFireLaser) {
      resp.push(["fireLaser"]);
    }

    if (shouldFireRocket && !shouldFireLaser) {
      resp.push(["fireRocket"]);
    }

    return resp;
  }
  onStart(spaceship: Spaceship, width: number, height: number) {
    this.startTs = Date.now();
  }
  onReset(spaceship: Spaceship, width: number, height: number) {
    this.startTs = Date.now();
  }
}

export const createKodyShip = () => async () => {
  const spaceship = new KodySpaceshipManager();
  return spaceship;
};
