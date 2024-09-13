import { isSpaceship } from "../../../../spaceships";
import type { GameObject } from "../../../../spaceships";


export type ScoreboardEntry = {
  id: number;
  name: string;
  score: number;
  kills: number;
  destroyed: boolean;
};

export const getScoreboard = (gameObjects: GameObject[]): ScoreboardEntry[] => {
  const states: ScoreboardEntry[] = [];
  for (const gameObject of gameObjects) {
    if (isSpaceship(gameObject)) {
      states.push({
        id: gameObject.id,
        name: gameObject.name,
        score: gameObject.score,
        kills: gameObject.kills,
        destroyed: gameObject.destroyed,
      });
    }
  }

  states.sort((a, b) => {
    if (a.destroyed && !b.destroyed) {
      return 1;
    } else if (!a.destroyed && b.destroyed) {
      return -1;
    }
    return b.score - a.score;
  });

  return states;
};
