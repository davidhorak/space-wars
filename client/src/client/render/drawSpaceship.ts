import {
  COLOR_COLLIDER,
  COLOR_ENERGY,
  COLOR_HEALTH,
  COLOR_TEXT,
  TEXT_INFO,
  TEXT_NAME,
} from ".";
import { rotateVector2, translateVector2, Vector2 } from "../utils";
import { Render } from "./render";
import { AnimatedSprite, Sprite } from "./Sprite";
import type { Spaceship } from "../../../../spaceships/types";

const MAIN_THRUST_SIZE = 24;
const SIDE_THRUST_SIZE = 18;
const THRUST_SPRITE_DURATION_MS = 100;

const thrustSize = (baseSize: number, thrust: number) => {
  return baseSize / 2 + ((baseSize / 2) * thrust) / 100;
};

export const drawSpaceship = ({
  render,
  spaceship,
  sprite,
  thrustSprite,
  showCollider,
  showHealth,
  showEnergy,
  showName,
  elapsedTimeMs,
}: {
  render: Render;
  spaceship: Spaceship;
  sprite: Sprite;
  thrustSprite: AnimatedSprite;
  showCollider: boolean;
  showHealth: boolean;
  showEnergy: boolean;
  showName: boolean;
  elapsedTimeMs: number;
}) => {
  render.drawSprite(
    sprite,
    spaceship.position.x,
    spaceship.position.y,
    spaceship.collider.radius * 2,
    spaceship.collider.radius * 2,
    spaceship.rotation + Math.PI / 2
  );

  let thrustPosition: Vector2;
  const thrustSpriteFrame = Math.floor(
    (elapsedTimeMs / THRUST_SPRITE_DURATION_MS) % thrustSprite.sprites.length
  );

  if (spaceship.engine.mainThrust > 0) {
    thrustPosition = translateVector2(
      { x: spaceship.position.x, y: spaceship.position.y },
      { x: 0, y: spaceship.collider.radius }
    );
    thrustPosition = rotateVector2(
      thrustPosition,
      spaceship.rotation + Math.PI / 2,
      { x: spaceship.position.x, y: spaceship.position.y }
    );
    const mainThrustSize = thrustSize(
      MAIN_THRUST_SIZE,
      spaceship.engine.mainThrust
    );

    render.drawSprite(
      thrustSprite.sprites[thrustSpriteFrame],
      thrustPosition.x,
      thrustPosition.y,
      mainThrustSize,
      mainThrustSize,
      spaceship.rotation + Math.PI / 2
    );
  }

  if (spaceship.engine.leftThrust > 0) {
    thrustPosition = translateVector2(
      { x: spaceship.position.x + 2, y: spaceship.position.y },
      { x: -spaceship.collider.radius, y: spaceship.collider.radius - 9 }
    );
    thrustPosition = rotateVector2(
      thrustPosition,
      spaceship.rotation + Math.PI / 2,
      { x: spaceship.position.x, y: spaceship.position.y }
    );

    const leftThrustSize = thrustSize(
      SIDE_THRUST_SIZE,
      spaceship.engine.leftThrust
    );

    render.drawSprite(
      thrustSprite.sprites[thrustSpriteFrame],
      thrustPosition.x,
      thrustPosition.y,
      leftThrustSize,
      leftThrustSize,
      spaceship.rotation + Math.PI
    );
  }

  if (spaceship.engine.rightThrust > 0) {
    thrustPosition = translateVector2(
      { x: spaceship.position.x - 2, y: spaceship.position.y },
      { x: spaceship.collider.radius, y: spaceship.collider.radius - 9 }
    );
    thrustPosition = rotateVector2(
      thrustPosition,
      spaceship.rotation + Math.PI / 2,
      { x: spaceship.position.x, y: spaceship.position.y }
    );

    const rightThrustSize = thrustSize(
      SIDE_THRUST_SIZE,
      spaceship.engine.rightThrust
    );

    render.drawSprite(
      thrustSprite.sprites[thrustSpriteFrame],
      thrustPosition.x,
      thrustPosition.y,
      rightThrustSize,
      rightThrustSize,
      spaceship.rotation - Math.PI * 2
    );
  }

  if (showName) {
    render.drawText(
      TEXT_NAME,
      COLOR_TEXT,
      spaceship.name,
      spaceship.position.x,
      spaceship.position.y - spaceship.collider.radius - 6,
      true
    );
  }

  let infoYOffset = 0;

  if (showHealth) {
    render.drawText(
      TEXT_INFO,
      COLOR_TEXT,
      "H",
      spaceship.position.x - 19,
      spaceship.position.y + spaceship.collider.radius + 12 + infoYOffset
    );
    render.drawRectFilled(
      COLOR_HEALTH,
      spaceship.position.x - 10,
      spaceship.position.y + spaceship.collider.radius + 6 + infoYOffset,
      20 * (spaceship.health / 100),
      6
    );
    render.drawRect(
      COLOR_HEALTH,
      1,
      spaceship.position.x - 10,
      spaceship.position.y + spaceship.collider.radius + 6 + infoYOffset,
      20,
      6
    );
    render.drawText(
      TEXT_INFO,
      COLOR_TEXT,
      Math.round(spaceship.health).toString(),
      spaceship.position.x + 12,
      spaceship.position.y + spaceship.collider.radius + 12 + infoYOffset
    );
    infoYOffset += 10;
  }

  if (showEnergy) {
    render.drawText(
      TEXT_INFO,
      COLOR_TEXT,
      "E",
      spaceship.position.x - 19,
      spaceship.position.y + spaceship.collider.radius + 12 + infoYOffset
    );
    render.drawRectFilled(
      COLOR_ENERGY,
      spaceship.position.x - 10,
      spaceship.position.y + spaceship.collider.radius + 6 + infoYOffset,
      20 * (spaceship.energy / 100),
      6
    );
    render.drawRect(
      COLOR_ENERGY,
      1,
      spaceship.position.x - 10,
      spaceship.position.y + spaceship.collider.radius + 6 + infoYOffset,
      20,
      6
    );
    render.drawText(
      TEXT_INFO,
      COLOR_TEXT,
      Math.round(spaceship.energy).toString(),
      spaceship.position.x + 12,
      spaceship.position.y + spaceship.collider.radius + 12 + infoYOffset
    );
  }

  if (showCollider) {
    render.drawCircle(
      COLOR_COLLIDER,
      1,
      spaceship.collider.position.x,
      spaceship.collider.position.y,
      spaceship.collider.radius
    );
  }
};
