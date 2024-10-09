import { Render } from "./render";
import { Sprite } from "./Sprite";
import { COLOR_COLLIDER } from ".";
import type { Asteroid } from "../../../../spaceships/types";

export const drawAsteroid = ({
  render,
  asteroid,
  sprite,
  elapsedTimeMs,
  showCollider,
}: {
  render: Render;
  asteroid: Asteroid;
  sprite: Sprite;
  elapsedTimeMs: number;
  showCollider: boolean;
}) => {
  render.drawSprite(
    sprite,
    asteroid.position.x,
    asteroid.position.y,
    asteroid.radius * 2,
    asteroid.radius * 2,
    asteroid.position.x + elapsedTimeMs / 4000 + asteroid.radius * 2
  );

  if (showCollider) {
    render.drawCircle(
      COLOR_COLLIDER,
      1,
      asteroid.collider.position.x,
      asteroid.collider.position.y,
      asteroid.collider.radius
    );
  }
};
