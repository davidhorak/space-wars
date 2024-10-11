import { rotateVector2, translateVector2 } from "../vector2";

describe("client / utils / vector2", () => {
  test("translateVector2", () => {
    expect(translateVector2({ x: 1, y: 2 }, { x: 3, y: 4 })).toEqual({
      x: 4,
      y: 6,
    });
  });

  describe("rotateVector2", () => {
    test.each([
      [{ x: 1, y: 2 }, 0, { x: 1, y: 2 }],
      [{ x: 1, y: 2 }, Math.PI / 2, { x: -2, y: 1 }],
    ])("rotateVector2(%s, %s) = %s", (vector, angle, expected) => {
      const result = rotateVector2(vector, angle);
      expect(result.x).toBeCloseTo(expected.x, 2);
      expect(result.y).toBeCloseTo(expected.y, 2);
    });

    test("origin", () => {
      const result = rotateVector2({ x: 1, y: 2 }, Math.PI / 2, { x: 3, y: 4 });
      expect(result.x).toBeCloseTo(5, 2);
      expect(result.y).toBeCloseTo(2, 2);
    });
  });
});
