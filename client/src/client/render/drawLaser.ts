import { COLOR_COLLIDER } from ".";
import type { Laser } from "../../../../spaceships/types";
import { Render } from "./render";
import { Sprite } from "./Sprite";

export const drawLaser = ({
  render,
  laser,
  sprite,
  scale,
  showCollider,
}: {
  render: Render;
  laser: Laser;
  sprite: Sprite;
  scale: number;
  showCollider: boolean;
}) => {
  render.drawSprite(
    sprite,
    laser.position.x,
    laser.position.y,
    sprite.width * scale,
    sprite.height * scale,
    laser.rotation + Math.PI / 2
  );

  if (showCollider) {
    render.drawRect(
      COLOR_COLLIDER,
      1,
      laser.collider.position.x,
      laser.collider.position.y,
      laser.collider.size.width,
      laser.collider.size.height,
      laser.rotation - Math.PI / 2
    );
  }
};
