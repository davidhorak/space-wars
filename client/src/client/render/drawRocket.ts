import { COLOR_COLLIDER } from ".";
import { Render } from "./render";
import { Sprite } from "./Sprite";

export const drawRocket = ({
  render,
  rocket,
  sprite,
  scale,
  showCollider,
}: {
  render: Render;
  rocket: Rocket;
  sprite: Sprite;
  scale: number;
  showCollider: boolean;
}) => {
  render.drawSprite(
    sprite,
    rocket.position.x,
    rocket.position.y,
    sprite.width * scale,
    sprite.height * scale,
    rocket.rotation + Math.PI / 2
  );

  if (showCollider) {
    render.drawCircle(
      COLOR_COLLIDER,
      1,
      rocket.collider.position.x,
      rocket.collider.position.y,
      rocket.collider.radius
    );
  }
};
