import { isUndefined, reverse, shuffle } from "lodash/fp";
import { isSetEngineThrustAction } from "../../../spaceships/spaceshipAction";
import { loop as createLoop } from "./loop";
import { observable } from "./observable";
import {
  render as createRender,
  sprites as createSprites,
  tiles as createTiles,
  drawAsteroid,
} from "./render";
import { drawBackground } from "./render/drawBackground";
import { drawExplosion } from "./render/drawExplosion";
import { drawLaser } from "./render/drawLaser";
import { drawRocket } from "./render/drawRocket";
import { drawSpaceship } from "./render/drawSpaceship";
import { getScoreboard } from "./utils";
import { getCanvas, getStartLocations } from "./utils";
import type { ScoreboardEntry } from "./utils";
import type { Vector2 } from "./utils";

import type {
  GameState,
  Spaceship,
  SpaceshipManager,
} from "../../../spaceships";
import spaceshipFactories, {
  isAsteroid,
  isExplosion,
  isLaser,
  isRocket,
  isSpaceship,
} from "../../../spaceships";

const TILE_SIZE = 48;

type EngineProps = {
  canvasId: string;
  width: number;
  height: number;
  fps: number;
};

export const engine = async ({ canvasId, width, height, fps }: EngineProps) => {
  if (spaceWars === null) {
    throw new Error("Space Wars kernel is not initialized.");
  }

  const canvas = getCanvas(canvasId);
  canvas.width = width;
  canvas.height = height;

  const context = canvas.getContext("2d");
  if (context === null) {
    throw new Error("failed to get the canvas context.");
  }

  const loop = createLoop({ fps });
  const render = createRender(canvas, context);
  const sprites = createSprites();

  let gameState: GameState | undefined;
  let elapsedTimeMs = 0;
  let logSize = 0;

  const onStateChanged = observable<GameState["status"]>();
  const onLogsChanged = observable<GameState["logs"]>();
  const onScoreboardChanged = observable<ScoreboardEntry[]>();
  const spaceships: SpaceshipManager[] = [];
  const spaceshipNameIndexMap = new Map<string, number>();

  // Options
  let showCollider = false;
  let showEnergy = false;
  let showHealth = false;
  let showNames = false;

  const start = () => {
    spaceWars.start();
    spaceships.forEach((spaceship) => spaceship.onStart(width, height));
    gameState = spaceWars.state();
    onStateChanged.broadcast("running");
    onLogsChanged.broadcast([]);
    onScoreboardChanged.broadcast(getScoreboard(gameState.gameObjects));
  };

  const reset = (startLocations: Vector2[]) => {
    elapsedTimeMs = 0;
    logSize = 0;
    startLocations = shuffle(startLocations);
    spaceships.forEach((spaceship, index) => {
      spaceWars.action(
        "setStartPosition",
        spaceship.name,
        startLocations[index].x,
        startLocations[index].y,
        Math.random() * Math.PI * 2
      );
    });
    spaceWars.reset();
    spaceWars.start();
    spaceships.forEach((spaceship) => spaceship.onReset());
    gameState = spaceWars.state();
    onStateChanged.broadcast("running");
    onLogsChanged.broadcast([]);
    onScoreboardChanged.broadcast(getScoreboard(gameState.gameObjects));
  };

  const onUpdate = (deltaTimeMs: number, forced: boolean = false) => {
    if (gameState?.status !== "running" && !forced) {
      return;
    }

    elapsedTimeMs += deltaTimeMs;

    try {
      spaceWars.tick(deltaTimeMs);
      gameState = spaceWars.state();

      for (const spaceship of spaceships) {
        let selfIndex = spaceshipNameIndexMap.get(spaceship.name);
        if (isUndefined(selfIndex)) {
          selfIndex = gameState.gameObjects.findIndex(
            (gameObject) =>
              isSpaceship(gameObject) && gameObject.name === spaceship.name
          );
          if (isUndefined(selfIndex)) {
            console.error(
              `spaceship ${spaceship.name} not found in the game state`
            );
            continue;
          }
          spaceshipNameIndexMap.set(spaceship.name, selfIndex);
        }

        const actions = spaceship.onUpdate({
          deltaTimeMs,
          gameObjects: gameState.gameObjects,
          spaceship: gameState.gameObjects[selfIndex] as Spaceship,
        });

        for (const action of actions) {
          if (isSetEngineThrustAction(action)) {
            spaceWars.action(
              action[0],
              spaceship.name,
              action[1],
              action[2],
              action[3]
            );
          } else {
            spaceWars.action(action[0], spaceship.name);
          }
        }
      }
    } catch (error) {
      console.error("failed to update", error);
    }
  };

  const onRender = (forced: boolean = false) => {
    if (!gameState || (gameState.status !== "running" && !forced)) {
      return;
    }

    try {
      if (gameState.logs.length > logSize) {
        onStateChanged.broadcast(gameState.status);
        onLogsChanged.broadcast(reverse(gameState.logs));
        onScoreboardChanged.broadcast(getScoreboard(gameState.gameObjects));
        logSize = gameState.logs.length;
      }

      drawBackground({
        render,
        backgroundPatterns,
        elapsedTimeMs,
        width,
        height,
      });

      for (let i = 0; i < gameState.gameObjects.length; i++) {
        const gameObject = gameState.gameObjects[i];
        if (!gameObject.enabled) {
          continue;
        }

        if (isAsteroid(gameObject)) {
          drawAsteroid({
            render,
            asteroid: gameObject,
            sprite: tilesMain.getTileByName(`asteroid_${i % 2}`),
            elapsedTimeMs: elapsedTimeMs * (i % 2 == 0 ? 1 : -1),
            showCollider,
          });
        } else if (isSpaceship(gameObject)) {
          drawSpaceship({
            render,
            spaceship: gameObject,
            sprite: tilesMain.getTileByName("spaceship"),
            thrustSprite: thrustSprite,
            elapsedTimeMs: elapsedTimeMs,
            showHealth,
            showEnergy,
            showName: showNames,
            showCollider,
          });
        } else if (isLaser(gameObject)) {
          drawLaser({
            render,
            laser: gameObject,
            sprite: tilesMain.getTileByName("laser"),
            scale: 0.75,
            showCollider,
          });
        } else if (isRocket(gameObject)) {
          drawRocket({
            render,
            rocket: gameObject,
            sprite: tilesMain.getTileByName("rocket"),
            scale: 0.75,
            showCollider,
          });
        } else if (isExplosion(gameObject)) {
          drawExplosion({
            render,
            explosion: gameObject,
            sprite: explosionsSprite,
          });
        }
      }
    } catch (error) {
      // TODO: Add popup with error message
      console.error("failed to render", error);
    }
  };

  try {
    await sprites.load(
      ["background_0", "img/space-wars_background_0.png", false],
      ["background_1", "img/space-wars_background_1.png", false],
      ["background_2", "img/space-wars_background_2.png", false],
      ["main", "img/space-wars_sprite.png", false]
    );
  } catch (error) {
    console.error("failed to load sprites", error);
  }

  const tilesMain = createTiles(sprites.get("main"), TILE_SIZE);
  tilesMain.mapTiles([
    ["asteroid_0", 0, 0],
    ["asteroid_1", 0, 1],
    ["spaceship", 1, 0],
    ["laser", 8, 3, 0.5],
    ["rocket", 8, 2, 0.5],
    ["explosion_0", 2, 0],
    ["explosion_1", 2, 1],
    ["explosion_2", 2, 2],
    ["explosion_3", 2, 3],
    ["explosion_4", 2, 4],
    ["explosion_5", 3, 0],
    ["explosion_6", 3, 1],
    ["explosion_7", 3, 2],
    ["explosion_8", 3, 3],
    ["explosion_9", 3, 4],
    ["thrust_0", 8, 0, 0.5],
    ["thrust_1", 8, 1, 0.5],
    ["thrust_2", 9, 0, 0.5],
    ["thrust_3", 9, 1, 0.5],
  ]);

  const backgroundPatterns = [
    render.createPattern(
      createTiles(sprites.get("background_0"), TILE_SIZE * 5).getTile(0, 0)
    ),
    render.createPattern(
      createTiles(sprites.get("background_1"), TILE_SIZE * 5).getTile(0, 0)
    ),
    render.createPattern(
      createTiles(sprites.get("background_2"), TILE_SIZE * 5).getTile(0, 0)
    ),
  ];

  const explosionsSprite = tilesMain.getAnimatedSpriteByName([
    "explosion_0",
    "explosion_1",
    "explosion_2",
    "explosion_3",
    "explosion_4",
    "explosion_5",
    "explosion_6",
    "explosion_7",
    "explosion_8",
    "explosion_9",
  ]);

  const thrustSprite = tilesMain.getAnimatedSpriteByName([
    "thrust_0",
    "thrust_1",
    "thrust_2",
    "thrust_3",
  ]);

  const startLocations = getStartLocations(
    width,
    height,
    spaceshipFactories.length
  );

  for (const [index, factory] of spaceshipFactories.entries()) {
    const instance = await factory();
    spaceWars.addSpaceship(
      instance.name,
      startLocations[index].x,
      startLocations[index].y,
      Math.random() * Math.PI * 2
    );
    spaceships.push(instance);
  }

  loop.onUpdate(onUpdate);
  loop.onRender(onRender);

  start();

  return {
    start: start,
    reset: () => reset(startLocations),
    pause: () => {
      spaceWars.pause();
      gameState = spaceWars.state();
      onStateChanged.broadcast("paused");
      onLogsChanged.broadcast(gameState.logs.reverse());
      onScoreboardChanged.broadcast(getScoreboard(gameState.gameObjects));
    },
    step: () => {
      onUpdate(50, true);
      onRender(true);
    },
    onStateChanged: (callback: (status: GameState["status"]) => void) =>
      onStateChanged.subscribe(callback),
    onLogsChanged: (callback: (logs: GameState["logs"]) => void) =>
      onLogsChanged.subscribe(callback),
    onScoreboardChanged: (callback: (scoreboard: ScoreboardEntry[]) => void) =>
      onScoreboardChanged.subscribe(callback),
    showCollider: (state: boolean) => (showCollider = state),
    showEnergy: (state: boolean) => (showEnergy = state),
    showHealth: (state: boolean) => (showHealth = state),
    showNames: (state: boolean) => (showNames = state),
  };
};

export type Engine = Awaited<ReturnType<typeof engine>>;
