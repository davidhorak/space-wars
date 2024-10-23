import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
  Spaceship,
} from "..";

const dummy = (name: string): SpaceshipManager => {
  let engineThrustEnergyTrigger = 100;
  let attacking = true;
  let lastRocketFireMs = 0;

  const onUpdate = ({
    spaceship: self,
    deltaTimeMs,
    gameObjects, // Add this line
  }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (lastRocketFireMs > 0) {
      lastRocketFireMs -= deltaTimeMs;
      if (lastRocketFireMs <= 0) {
        lastRocketFireMs = 0;
      }
    }

    for (const gameObject of gameObjects) {
      if (gameObject.type === "spaceship") {
        const distance = Math.sqrt(
          Math.pow(gameObject.position.x - self.position.x, 2) +
            Math.pow(gameObject.position.y - self.position.y, 2)
        );

        if (distance < 400) {
          if (self.energy < 5) {
            actions.push(["setEngineThrust", 0, 0, 0]);
            continue;
          }

          if (self.energy > 30) {
            actions.push(["fireRocket"]);
          }

          actions.push(["fireLaser"]);

          const angleToTarget = Math.atan2(
            gameObject.position.y - self.position.y,
            gameObject.position.x - self.position.x
          );

          const angleDifference = angleToTarget - self.rotation;

          if (Math.abs(angleDifference) < Math.PI / 2) {
            if (angleDifference > 0) {
              actions.push(["setEngineThrust", 10, 0, 0]);
            } else {
              actions.push(["setEngineThrust", 10, 0, 0]);
            }
          }
        }
      } else if (gameObject.type === "asteroid") {
        const distance = Math.sqrt(
          Math.pow(gameObject.position.x - self.position.x, 2) +
            Math.pow(gameObject.position.y - self.position.y, 2)
        );

        if (distance < 150) {
          const angleToAsteroid = Math.atan2(
            gameObject.position.y - self.position.y,
            gameObject.position.x - self.position.x
          );

          const angleDifference = angleToAsteroid - self.rotation;

          if (Math.abs(angleDifference) < Math.PI / 2) {
            actions.push(["setEngineThrust", 10, 0, 0]);
          }
        }
      }
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: function (spaceship, width, height) {
      this.onReset(spaceship, width, height);
    },
    onReset: () => {
      attacking = true;
      lastRocketFireMs = 0;
      engineThrustEnergyTrigger = 100;
    },
  };
};

export const createLethalrush = (): SpaceshipManagerFactory => () =>
  // Promise.resolve(dummy(names[++index]));
  Promise.resolve(dummy("lethalrush"));
