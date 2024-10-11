import { SpaceshipAction } from "./spaceshipAction";
import { GameState, Spaceship } from "./types";

export type SpaceState = {
  deltaTimeMs: number;
  spaceship: Spaceship;
  gameObjects: GameState["gameObjects"];
};

export interface SpaceshipManager {
  name: string;
  onUpdate(state: SpaceState): SpaceshipAction[];
  onStart(spaceship: Spaceship, width: number, height: number): void;
  onReset(spaceship: Spaceship, width: number, height: number): void;
}

export type SpaceshipManagerFactory = () => Promise<SpaceshipManager>;
