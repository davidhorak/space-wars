type MainEngineThrust = number;
type LeftEngineThrust = number;
type RightEngineThrust = number;

export type SetEngineThrustAction = [
  "setEngineThrust",
  MainEngineThrust,
  LeftEngineThrust,
  RightEngineThrust
];

type X = number;
type Y = number;
type Rotation = number;

export type SetStartPositionAction = ["setStartPosition", X, Y, Rotation];

export type FireLaserAction = ["fireLaser"];
export type FireRocketAction = ["fireRocket"];

export type SpaceshipAction =
  | SetEngineThrustAction
  | FireLaserAction
  | FireRocketAction;

export const isSetEngineThrustAction = (
  action: unknown[]
): action is SetEngineThrustAction => {
  return action[0] === "setEngineThrust";
};

export const isSetStartPositionAction = (
  action: unknown[]
): action is SetStartPositionAction => {
  return action[0] === "setStartPosition";
};

export const isFireLaserAction = (
  action: unknown[]
): action is FireLaserAction => {
  return action[0] === "fireLaser";
};

export const isFireRocketAction = (
  action: unknown[]
): action is FireRocketAction => {
  return action[0] === "fireRocket";
};
