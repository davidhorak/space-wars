import { getStartLocations } from "../startLocations";

describe("client / utils / getStartLocations", () => {
  it.each([
    [10, 10, 2, [{ x: 8.333, y: 5 }, { x: 1.666, y: 5 }]],
    [10, 10, 4, [{ x: 8.333, y: 5 }, { x: 5, y: 8.333 }, { x: 1.666, y: 5 }, { x: 5, y: 1.666 }]],
  ])("returns %s slots for a %s x %s screen", (width, height, slots, expected) => {
    const locations = getStartLocations(width, height, slots);
    expect(locations).toHaveLength(slots);

    expected.forEach((location, index) => {
      expect(locations[index].x).toBeCloseTo(location.x, 2);
      expect(locations[index].y).toBeCloseTo(location.y, 2);
    });
  });
});

