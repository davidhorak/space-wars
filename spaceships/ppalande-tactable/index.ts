import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
  GameObject,
  Spaceship,
  Laser,
  Rocket,
} from "..";

const ppalandeTactable = (name: string): SpaceshipManager => {
  let lastRocketFireMs = 0;
  let isRecharging = false;

  const MIN_ENERGY_FOR_ACTIONS = 30;
  const RECHARGE_THRESHOLD = 70;


  const findNearestSpaceship = (self: Spaceship, gameObjects: GameObject[]): Spaceship | null => {
    let nearestDistance = Infinity;
    let nearestSpaceship: Spaceship | null = null;

    gameObjects.forEach(object => {
      if (object.type === "spaceship" && object.id !== self.id) {
        const distance = Math.hypot(object.position.x - self.position.x, object.position.y - self.position.y);
        if (distance < nearestDistance) {
          nearestDistance = distance;
          nearestSpaceship = object as Spaceship;
        }
      }
    });

    return nearestSpaceship;
  };

  const findNearestAsteroid = (self: Spaceship, gameObjects: GameObject[]): GameObject | null => {
    let nearestDistance = Infinity;
    let nearestAsteroid: GameObject | null = null;

    gameObjects.forEach(object => {
      if (object.type === "asteroid") {
        const distance = Math.hypot(object.position.x - self.position.x, object.position.y - self.position.y);
        if (distance < nearestDistance) {
          nearestDistance = distance;
          nearestAsteroid = object;
        }
      }
    });

    const avoidAngle = nearestAsteroid.position.angle + Math.PI; // Turn 180 degrees to avoid the asteroid
    return nearestAsteroid;
  };

  const findIncomingProjectile = (self: Spaceship, gameObjects: GameObject[]): GameObject | null => {
    let nearestProjectile: GameObject | null = null;
    let nearestDistance = Infinity;

    gameObjects.forEach(object => {
      if ((object.type === "laser" || object.type === "rocket") && object.owner.id !== self.id) {
        const distance = Math.hypot(object.position.x - self.position.x, object.position.y - self.position.y);
        if (distance < nearestDistance) {
          nearestDistance = distance;
          nearestProjectile = object;
        }
      }
    });

    return nearestProjectile;
  };

  const calculateAngleToTarget = (self: Spaceship, target: GameObject): number => {
    const dx = target.position.x - self.position.x;
    const dy = target.position.y - self.position.y;
    return Math.atan2(dy, dx);
  };

  const onUpdate = ({
    spaceship: self,
    deltaTimeMs,
    gameObjects,
  }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) return actions;

    if (lastRocketFireMs > 0) {
      lastRocketFireMs -= deltaTimeMs;
      if (lastRocketFireMs <= 0) lastRocketFireMs = 0;
    }

    // Energy management
    if (self.energy <= MIN_ENERGY_FOR_ACTIONS) {
      isRecharging = true;
    } else if (self.energy >= RECHARGE_THRESHOLD) {
      isRecharging = false;
    }

    if (isRecharging) {
      return actions; // Return empty actions to recharge
    }

    const nearestSpaceship = findNearestSpaceship(self, gameObjects);
    const incomingProjectile = findIncomingProjectile(self, gameObjects);
    const nearestAsteroid = findNearestAsteroid(self, gameObjects);

    let targetAngle: number;

    // 1. Avoid incoming projectiles
    if (incomingProjectile) {
      targetAngle = calculateAngleToTarget(self, incomingProjectile) + Math.PI / 2;
    }
    // 2. Avoid incoming spaceships
    else if (nearestSpaceship && Math.hypot(nearestSpaceship.position.x - self.position.x, nearestSpaceship.position.y - self.position.y) < 100) {
      targetAngle = calculateAngleToTarget(self, nearestSpaceship) + Math.PI;
    }
    // 3. Avoid running into asteroids
    else if (nearestAsteroid && Math.hypot(nearestAsteroid.position.x - self.position.x, nearestAsteroid.position.y - self.position.y) < 100) {
      targetAngle = calculateAngleToTarget(self, nearestAsteroid) + Math.PI;
    }
    // 4. Point towards the closest spaceship
    else if (nearestSpaceship) {
      targetAngle = calculateAngleToTarget(self, nearestSpaceship);
    }
    else {
      // No threats or targets, turn randomly
      targetAngle = Math.random() * Math.PI * 2;
    }

    const angleDiff = (targetAngle - self.rotation + Math.PI * 3) % (Math.PI * 2) - Math.PI;

    let mainThrust = 100;
    let leftThrust = 0;
    let rightThrust = 0;

    if (Math.abs(angleDiff) > 0.1) {
      mainThrust = 20;
      if (angleDiff > 0) {
        leftThrust = 80;
      } else {
        rightThrust = 80;
      }
    }

    actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);

    // Shoot at the nearest spaceship
    if (nearestSpaceship && Math.abs(angleDiff) < 0.1) {
      const distance = Math.hypot(nearestSpaceship.position.x - self.position.x, nearestSpaceship.position.y - self.position.y);
      
      if (Math.random() < 0.5) {  // 50% probability to choose between laser and rocket
        if (distance < 200 && self.rocketReloadTimerSec === 0 && self.rockets > 0 && self.energy >= 20) {
          actions.push(["fireRocket"]);
          lastRocketFireMs = 500;
        }
      } else {
        if (distance < 150 && self.laserReloadTimerSec === 0 && self.energy >= 10) {
          actions.push(["fireLaser"]);
        }
      }
    }

    return actions;
  };

  return {
    name,
    onUpdate,
    onStart: () => {},
    onReset: () => {
      lastRocketFireMs = 0;
      isRecharging = false;
    },
  };
};

export const createPpalandeTactable = (): SpaceshipManagerFactory => () =>
  Promise.resolve(ppalandeTactable("dattebayo"));
