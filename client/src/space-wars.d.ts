declare namespace spaceWars {
  function init(
    width: number,
    height: number,
    seed?: number
  ): void;
  function tick(deltaTimeMs: number): void;
  function start(): void;
  function pause(): void;
  function reset(): void;
  function state(): import('../../spaceships').GameState;
  function fromState(state: string): void;
  function addSpaceship(
    name: string,
    x: number,
    y: number,
    rotation: number
  ): void;
  function action(
    action: "setEngineThrust",
    shipName: string,
    mainEngineThrust: number,
    leftEngineThrust: number,
    rightEngineThrust: number
  ): void;
  function action(action: "fireLaser" | "fireRocket", shipName: string): void;
  function action(action: "setStartPosition", shipName: string, x: number, y: number, rotation: number): void;
}
