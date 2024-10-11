import { isUndefined } from "lodash/fp";
import { GameObject, isSpaceship, Spaceship } from "../../../../spaceships";

export const getSpaceship = (
  name: string,
  gameObjects: GameObject[],
  spaceshipNameIndexMap: Map<string, number>
): Spaceship => {
  let selfIndex = spaceshipNameIndexMap.get(name);
  if (isUndefined(selfIndex)) {
    selfIndex = gameObjects.findIndex(
      (gameObject) => isSpaceship(gameObject) && gameObject.name === name
    );
    if (selfIndex < 0) {
      throw new Error(`spaceship ${name} not found in the game state`);
    }
    spaceshipNameIndexMap.set(name, selfIndex);
  }
  return gameObjects[selfIndex] as Spaceship;
};
