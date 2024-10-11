import { GameObject, Spaceship } from "../../../../../spaceships";
import { getScoreboard } from "../scoreboard";

test("client / utils / getScoreboard", () => {
  const gameObjects: GameObject[] = [
    {
      id: 1,
      name: "Spaceship 1 (destroyed)",
      score: 100,
      kills: 0,
      destroyed: true,
      type: "spaceship",
    } as Spaceship,
    {
      id: 2,
      name: "Spaceship 2",
      score: 100,
      kills: 0,
      destroyed: false,
      type: "spaceship",
    } as Spaceship,
    {
      id: 3,
      name: "Spaceship 3",
      score: 200,
      kills: 3,
      destroyed: false,
      type: "spaceship",
    } as Spaceship,
    {
      id: 4,
      name: "Spaceship 4 (destroyed)",
      score: 200,
      kills: 0,
      destroyed: true,
      type: "spaceship",
    } as Spaceship,
  ];

  const scoreboard = getScoreboard(gameObjects)

  expect(scoreboard.map((entry) => entry.name)).toEqual([
    "Spaceship 3",
    "Spaceship 2",
    "Spaceship 4 (destroyed)",
    "Spaceship 1 (destroyed)",
  ]);
  expect(scoreboard[0].score).toBe(200);
  expect(scoreboard[0].kills).toBe(3);
  expect(scoreboard[0].destroyed).toBe(false);
});
