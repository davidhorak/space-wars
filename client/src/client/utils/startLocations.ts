import { Vector2 } from "./vector2";

// Evenly distributes slots around the perimeter of the screen on circle
export const getStartLocations = (
  width: number,
  height: number,
  slots: number
): Vector2 [] => {
  const startLocations: Vector2[] = [];
  const center = { x: width / 2, y: height / 2 };
  const radius = Math.min(width, height) / 3;

  for (let i = 0; i < slots; i++) {
    const angle = (i / slots) * 2 * Math.PI;
    const startX = center.x + radius * Math.cos(angle);
    const startY = center.y + radius * Math.sin(angle);
    startLocations.push({ x: startX, y: startY });
  }

  return startLocations;
};
