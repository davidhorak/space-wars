import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
} from "./spaceshipManager";
import { createDummy } from "./_samples/dummy";
import { createPpalandeTactable } from "./ppalande-tactable";
import { createRocinante } from "./gordon-tactable/rocinante";
import { initShip as initVicecarloans } from "./vicecarloans";
import { createAmazingLilySpaceship } from "./lily";
import { createJustin } from "./justin";
import { createKodyShip } from "./kodyang";
import { createTnT } from "./mason";
import { createLethalrush } from "./saro";

// import {
//   createEngineTester,
//   createLaserTester,
//   createRocketTester,
// } from "./_samples/tester";

export {
  isAsteroid,
  isExplosion,
  isLaser,
  isRocket,
  isSpaceship,
} from "./gameObject";

export type { SpaceshipAction } from "./spaceshipAction";
export type { SpaceshipManagerFactory, SpaceshipManager, SpaceState };
export type {
  SquareCollider,
  CircleCollider,
  PolygonCollider,
  GameObject,
  Spaceship,
  Log,
  DamageLog,
  KillLog,
  CollisionLog,
  GameStateLog,
  GameState,
} from "./types";

const spaceshipFactories: SpaceshipManagerFactory[] = [
  // createDummy(),
  // createDummy(),
  // createDummy(),
  createDummy(),
  createJustin(),
  createRocinante("Rocinante"),
  initVicecarloans(),
  createAmazingLilySpaceship(),
  createPpalandeTactable(),
  createTnT(),
  createKodyShip(),
  createLethalrush(),
  //   createLaserTester("Laser Tester 1"),
  //   createLaserTester("Laser Tester 2"),
  //   createRocketTester("Rocket Tester 1"),
  //   createRocketTester("Rocket Tester 2"),
  // createEngineTester("Engine Tester MAIN", 100, 0, 0),
  // createEngineTester("Engine Tester LEFT", 25, 100, 0),
  // createEngineTester("Engine Tester RIGHT", 25, 0, 100),
];

export default spaceshipFactories;
