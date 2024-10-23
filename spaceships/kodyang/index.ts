import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
  Spaceship,
} from "..";

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

    const resp: SpaceshipAction[] = [
      ["setEngineThrust", 0, state.spaceship.energy <= 20 ? 0 : 0.1, 0],
    ];

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
