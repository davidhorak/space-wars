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
  const MIN_ENERGY = 30;
  const MAX_ENERGY = 70;
  const DODGE_DISTANCE = 100;
  const ASTEROID_AVOID_DISTANCE = 50;
  const LASER_RANGE = 150;
  const ROCKET_RANGE = 200;

  const findNearestObject = (self: Spaceship, gameObjects: GameObject[], type: string): GameObject | null => {
    let nearest: GameObject | null = null;
    let minDistance = Infinity;
    gameObjects.forEach(obj => {
      if (obj.type === type && obj.id !== self.id) {
        const distance = Math.hypot(obj.position.x - self.position.x, obj.position.y - self.position.y);
        if (distance < minDistance) {
          minDistance = distance;
          nearest = obj;
        }
      }
    });
    return nearest;
  };

  const calculateAngleToTarget = (self: Spaceship, target: GameObject): number => {
    const dx = target.position.x - self.position.x;
    const dy = target.position.y - self.position.y;
    return Math.atan2(dy, dx);
  };

  const predictPosition = (obj: GameObject, timeAhead: number): { x: number, y: number } => {
    return {
      x: obj.position.x + obj.velocity.x * timeAhead,
      y: obj.position.y + obj.velocity.y * timeAhead
    };
  };

  const shouldDodge = (self: Spaceship, gameObjects: GameObject[]): GameObject | null => {
    return gameObjects.find(obj => 
      (obj.type === "laser" || obj.type === "rocket") && 
      obj.owner.id !== self.id &&
      Math.hypot(obj.position.x - self.position.x, obj.position.y - self.position.y) < DODGE_DISTANCE
    ) || null;
  };

  const avoidAsteroids = (self: Spaceship, gameObjects: GameObject[]): { x: number, y: number } => {
    let avoidVector = { x: 0, y: 0 };
    gameObjects.forEach(obj => {
      if (obj.type === "asteroid") {
        const distance = Math.hypot(obj.position.x - self.position.x, obj.position.y - self.position.y);
        if (distance < ASTEROID_AVOID_DISTANCE) {
          const angle = Math.atan2(self.position.y - obj.position.y, self.position.x - obj.position.x);
          avoidVector.x += Math.cos(angle) * (ASTEROID_AVOID_DISTANCE - distance);
          avoidVector.y += Math.sin(angle) * (ASTEROID_AVOID_DISTANCE - distance);
        }
      }
    });
    return avoidVector;
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

    const nearestEnemy = findNearestObject(self, gameObjects, "spaceship") as Spaceship;
    const dodgeObject = shouldDodge(self, gameObjects);
    const asteroidAvoidance = avoidAsteroids(self, gameObjects);

    let targetAngle: number;
    let mainThrust = 0;
    let leftThrust = 0;
    let rightThrust = 0;

    // Dodge incoming projectiles or avoid asteroids
    if (dodgeObject || (asteroidAvoidance.x !== 0 || asteroidAvoidance.y !== 0)) {
      const avoidAngle = dodgeObject ? 
        calculateAngleToTarget(self, dodgeObject) + Math.PI / 2 :
        Math.atan2(asteroidAvoidance.y, asteroidAvoidance.x);
      targetAngle = avoidAngle;
      mainThrust = 100;
    } 
    // Attack nearest enemy
    else if (nearestEnemy && self.energy > MIN_ENERGY) {
      const distanceToEnemy = Math.hypot(nearestEnemy.position.x - self.position.x, nearestEnemy.position.y - self.position.y);
      const predictedPosition = predictPosition(nearestEnemy, distanceToEnemy / 320); // 320 is laser speed
      targetAngle = calculateAngleToTarget(self, { position: predictedPosition } as GameObject);
      
      const angleDiff = (targetAngle - self.rotation + Math.PI * 3) % (Math.PI * 2) - Math.PI;

      if (Math.abs(angleDiff) < 0.1) {
        if (distanceToEnemy <= ROCKET_RANGE && self.rockets > 0 && self.rocketReloadTimerSec === 0 && self.energy >= 20) {
          actions.push(["fireRocket"]);
          lastRocketFireMs = 1000;
        } else if (distanceToEnemy <= LASER_RANGE && self.laserReloadTimerSec === 0 && self.energy >= 6) {
          actions.push(["fireLaser"]);
        }
      }

      mainThrust = 50; // Maintain some movement for unpredictability
    } 
    // Recharge energy
    else {
      targetAngle = self.rotation; // Maintain current rotation
      mainThrust = 0; // Stop to recharge faster
    }

    // Apply rotation
    const angleDiff = (targetAngle - self.rotation + Math.PI * 3) % (Math.PI * 2) - Math.PI;
    if (Math.abs(angleDiff) > 0.1) {
      if (angleDiff > 0) {
        leftThrust = 80;
      } else {
        rightThrust = 80;
      }
    }

    // Energy management
    if (self.energy > MAX_ENERGY) {
      mainThrust = Math.max(mainThrust, 20); // Use some energy if we have excess
    } else if (self.energy < MIN_ENERGY) {
      mainThrust = 0; // Stop to recharge
      leftThrust = 0;
      rightThrust = 0;
    }

    actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);

    return actions;
  };

  return {
    name,
    onUpdate,
    onStart: () => {},
    onReset: () => { lastRocketFireMs = 0; },
  };
};

export const createPpalandeTactable = (): SpaceshipManagerFactory => () =>
  Promise.resolve(ppalandeTactable("dattebayo"));
