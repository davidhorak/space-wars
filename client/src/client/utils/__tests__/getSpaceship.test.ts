import { Spaceship } from "../../../../../spaceships";
import { getSpaceship } from "../getSpaceship";

describe("client / utils / getSpaceship", () => {
  test("found in the map", () => {
    const name = "Deathwatch";
    const gameObjects = [{ id: 1 } as Spaceship];
    const spaceshipNameIndexMap = new Map([[name, 0]]);

    expect(getSpaceship(name, gameObjects, spaceshipNameIndexMap).id).toBe(1);
    expect(spaceshipNameIndexMap.size).toBe(1);
  });

  test("not found in the map", () => {
    const name = "Vindicator";
    const gameObjects = [{ id: 1, name: "Vindicator", type: "spaceship" } as Spaceship];
    const spaceshipNameIndexMap = new Map();

    expect(getSpaceship(name, gameObjects, spaceshipNameIndexMap).id).toBe(1);
    expect(spaceshipNameIndexMap.size).toBe(1);
  });

  test("not found in the game objects", () => {
    expect(() => getSpaceship("Vindicator", [], new Map())).toThrow(
      'spaceship Vindicator not found in the game state'
    );
  });
});
