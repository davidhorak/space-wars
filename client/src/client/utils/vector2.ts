export interface Vector2 {
  x: number;
  y: number;
}

export const translateVector2 = (vector: Vector2, translation: Vector2): Vector2 => {
  return {
    x: vector.x + translation.x,
    y: vector.y + translation.y,
  };
};

export const rotateVector2 = (vector: Vector2, angle: number, origin: Vector2 = { x: 0, y: 0 }): Vector2 => {
  const cos = Math.cos(angle);
  const sin = Math.sin(angle);

  const x = vector.x - origin.x;
  const y = vector.y - origin.y;

  return {
    x: x * cos - y * sin + origin.x,
    y: x * sin + y * cos + origin.y,
  };
};