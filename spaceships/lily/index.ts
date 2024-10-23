import { SpaceshipAction } from "../spaceshipAction";
import {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
} from "../spaceshipManager";

const MY_NAME = "POTATO!";

function isTargetInAim(
  x: number,
  y: number,
  rotation: number,
  targetX: number,
  targetY: number,
  aimRadius: number,
  maxDistance: number,
  gameWidth: number,
  gameHeight: number
): boolean {
  const deltaX = ((targetX - x + gameWidth / 2) % gameWidth) - gameWidth / 2;
  const deltaY = ((targetY - y + gameHeight / 2) % gameHeight) - gameHeight / 2;

  const aimDirX = Math.cos(rotation);
  const aimDirY = Math.sin(rotation);

  const distanceToTarget = Math.sqrt(deltaX ** 2 + deltaY ** 2);

  if (distanceToTarget > maxDistance) {
    return false;
  }

  const projectionLength = deltaX * aimDirX + deltaY * aimDirY;
  const perpendicularDist = Math.abs(deltaX * aimDirY - deltaY * aimDirX);
  return projectionLength >= 0 && perpendicularDist <= aimRadius;
}

function getTurnDirection(
  x: number,
  y: number,
  rotation: number,
  targetX: number,
  targetY: number,
  gameWidth: number,
  gameHeight: number
): "left" | "right" | "none" {
  const deltaX = ((targetX - x + gameWidth / 2) % gameWidth) - gameWidth / 2;
  const deltaY = ((targetY - y + gameHeight / 2) % gameHeight) - gameHeight / 2;

  const aimDirX = Math.cos(rotation);
  const aimDirY = Math.sin(rotation);

  // Calculate the cross product to determine turn direction
  const crossProduct = deltaX * aimDirY - deltaY * aimDirX;

  if (crossProduct > 0) {
    return "left";
  } else if (crossProduct < 0) {
    return "right";
  } else {
    return "none";
  }
}

export const AmazingSpaceship = (): SpaceshipManager => {
  let time = 0;
  let thrustTimer = 0;
  let thurstEndTimer = 500;
  let gameWidth = 0;
  let gameHeight = 0;
  let currentThrust = 0;
  const onUpdate = ({
    spaceship: self,
    deltaTimeMs,
    gameObjects,
  }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];

    if (thrustTimer <= 0) {
      thrustTimer = 1000;

      actions.push(["setEngineThrust", 100, 0, 0]);
      currentThrust = 100;
    }

    if (thurstEndTimer <= 0) {
      thurstEndTimer = 1000;
      actions.push(["setEngineThrust", 0, 0, 0]);
      currentThrust = 0;
    }

    thrustTimer -= deltaTimeMs;
    thurstEndTimer -= deltaTimeMs;

    for (const gameObject of gameObjects) {
      if (gameObject.type !== "spaceship" || gameObject.name === MY_NAME) {
        continue;
      }

      if (
        isTargetInAim(
          self.position.x,
          self.position.y,
          self.rotation,
          gameObject.position.x,
          gameObject.position.y,
          gameObject.collider.radius,
          1000,
          gameWidth,
          gameHeight
        )
      ) {
        actions.push(["fireLaser"]);
      }

      if (
        isTargetInAim(
          self.position.x,
          self.position.y,
          self.rotation,
          gameObject.position.x,
          gameObject.position.y,
          gameObject.collider.radius * 20,
          1000,
          gameWidth,
          gameHeight
        )
      ) {
        const direction = getTurnDirection(
          self.position.x,
          self.position.y,
          self.rotation,
          gameObject.position.x,
          gameObject.position.y,
          gameWidth,
          gameHeight
        );

        // if (direction == "left") {
        //   actions.push(["setEngineThrust", currentThrust, 0, 40]);
        // }
        // if (direction == "right") {
        //   actions.push(["setEngineThrust", currentThrust, 40, 0]);
        // }
      }
      gameObject.position.x;
    }

    for (const gameObject of gameObjects) {
      if (gameObject.type !== "asteroid") {
        continue;
      }

      if (
        isTargetInAim(
          self.position.x,
          self.position.y,
          self.rotation,
          gameObject.position.x,
          gameObject.position.y,
          (gameObject.collider.radius + self.collider.radius) * 1.5,
          2000,
          gameWidth,
          gameHeight
        )
      ) {
        actions.push(["setEngineThrust", 20, 40, 0]);
      }
    }

    return actions;
  };

  return {
    name: MY_NAME,
    onUpdate: onUpdate,
    onStart: function (spaceship, width, height) {
      this.onReset(spaceship, width, height);
    },
    onReset: (spaceship, width, height) => {
      time = 0;
      gameWidth = width;
      gameHeight = height;
    },
  };
};

export const createAmazingLilySpaceship = (): SpaceshipManagerFactory => () =>
  Promise.resolve(AmazingSpaceship());