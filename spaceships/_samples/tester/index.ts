import {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
} from "../..";

const lasterTester = (name: string): SpaceshipManager => {
  const onUpdate = ({ spaceship: self }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (self.energy > 10 && self.laserReloadTimerSec == 0) {
      actions.push(["fireLaser"]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

const rocketTester = (name: string): SpaceshipManager => {
  const onUpdate = ({ spaceship: self }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (self.energy > 20 && self.rocketReloadTimerSec == 0) {
      actions.push(["fireRocket"]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

const engineTester = (
  name: string,
  mainEngineThrust: number,
  leftEngineThrust: number,
  rightEngineThrust: number
): SpaceshipManager => {
  const onUpdate = ({ spaceship: self }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (self.energy >= 100) {
      actions.push([
        "setEngineThrust",
        mainEngineThrust,
        leftEngineThrust,
        rightEngineThrust,
      ]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

export const createLaserTester =
  (name: string): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(lasterTester(name));

export const createRocketTester =
  (name: string): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(rocketTester(name));

export const createEngineTester =
  (
    name: string,
    mainEngineThrust: number,
    leftEngineThrust: number,
    rightEngineThrust: number
  ): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(
      engineTester(name, mainEngineThrust, leftEngineThrust, rightEngineThrust)
    );
