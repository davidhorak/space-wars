import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
} from "../..";

const dummy = (name: string): SpaceshipManager => {
  let engineThrustEnergyTrigger = 100;
  let attacking = true;
  let lastRocketFireMs = 0;

  const onUpdate = ({
    spaceship: self,
    deltaTimeMs,
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

    if (self.energy <= 10) {
      attacking = false;
      engineThrustEnergyTrigger = 25 + Math.random() * 50;
    }

    if (self.energy >= engineThrustEnergyTrigger) {
      engineThrustEnergyTrigger = Infinity;
      attacking = true;
      const mainThrust = 50 + Math.random() * 50;
      let leftThrust = 0;
      let rightThrust = 0;
      if (Math.random() < 0.5) {
        leftThrust = 10 + Math.random() * 90;
      } else {
        rightThrust = 10 + Math.random() * 90;
      }
      actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);
    }

    if (
      attacking &&
      self.energy > 20 &&
      self.rocketReloadTimerSec == 0 &&
      self.rockets > 0
    ) {
      actions.push(["fireRocket"]);
      lastRocketFireMs = 500;
    } else if (
      lastRocketFireMs <= 0 &&
      attacking &&
      self.energy > 10 &&
      self.laserReloadTimerSec == 0
    ) {
      actions.push(["fireLaser"]);
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

const names = [
  "Ultramar",
  "Dark Angel",
  "Blood Angel",
  "Space Wolf",
  "Imperial Fist",
  "White Scar",
  "Salamander",
  "Raven Guard",
  "Iron Hand",
  "Deathwatch",
];

let index = -1;

export const createDummy = (): SpaceshipManagerFactory => () =>
  Promise.resolve(dummy(names[++index]));
