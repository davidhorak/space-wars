import type { Explosion } from "../../../../spaceships/types";
import { Render } from "./render";
import { AnimatedSprite } from "./Sprite";

export const drawExplosion = ({
  render,
  explosion,
  sprite,
}: {
  render: Render;
  explosion: Explosion;
  sprite: AnimatedSprite;
}) => {
  const frameDuration = explosion.durationSec / sprite.sprites.length;
  const frame =
    sprite.sprites.length -
    1 -
    Math.floor(explosion.lifespanSec / frameDuration);

  render.drawSprite(
    sprite.sprites[frame],
    explosion.position.x,
    explosion.position.y,
    explosion.radius * 2,
    explosion.radius * 2
  );
};
