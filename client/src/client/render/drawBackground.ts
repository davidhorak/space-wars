import { Render } from "../render";
import { rotateVector2 } from "../utils/vector2";

export const drawBackground = ({
  render,
  backgroundPatterns,
  elapsedTimeMs,
  width,
  height,
}: {
  render: Render;
  backgroundPatterns: CanvasPattern[];
  elapsedTimeMs: number;
  width: number;
  height: number;
}) => {
  let angle = elapsedTimeMs * 0.001;
  let position = rotateVector2({ x: 0, y: 10 }, angle);

  render.translate(position.x, position.y);
  render.drawPattern(backgroundPatterns[0], -10, -10, width + 20, height + 20);
  render.translate(-position.x, -position.y);

  angle = (elapsedTimeMs - 250) * 0.0005;
  position = rotateVector2({ x: 0, y: 10 }, -angle);

  render.translate(0, position.y);
  render.drawPattern(backgroundPatterns[1], -10, -10, width + 20, height + 20);
  render.translate(0, -position.y);

  angle = (elapsedTimeMs - 500) * 0.0005;
  position = rotateVector2({ x: 0, y: 10 }, angle);

  render.translate(position.x, 0);
  render.drawPattern(backgroundPatterns[2], -10, -10, width + 20, height + 20);
  render.translate(-position.x, 0);
}