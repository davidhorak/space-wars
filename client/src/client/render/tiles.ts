import type { AnimatedSprite, Sprite } from "./Sprite";

const fallbackSprite: Sprite = {
  image: new Image(),
  x: 0,
  y: 0,
  width: 0,
  height: 0,
};

export const tiles = (image: HTMLImageElement, tileSize: number) => ({
  tileSize,
  tileMap: new Map<string, Sprite>(),
  mapTiles: function (
    mapping: [name: string, row: number, column: number, tileScale?: number][]
  ) {
    mapping.forEach(([name, row, column, tileScale]) => {
      this.tileMap.set(name, this.getTile(row, column, tileScale));
    });
  },
  getTile: (row: number, column: number, tileScale = 1): Sprite => ({
    image,
    x: column * (tileSize * tileScale),
    y: row * (tileSize * tileScale),
    width: tileSize * tileScale,
    height: tileSize * tileScale,
  }),
  getTileByName: function (name: string): Sprite {
    if (this.tileMap.has(name)) {
      return this.tileMap.get(name)!;
    }

    console.warn(`Tile not found: ${name}`);
    return fallbackSprite;
  },
  getAnimatedSprite: (
    rowColumnPairs: [number, number][],
    tileScale = 1
  ): AnimatedSprite => ({
    sprites: rowColumnPairs.map(([row, column]) => ({
      image,
      x: column * (tileSize * tileScale),
      y: row * (tileSize * tileScale),
      width: tileSize * tileScale,
      height: tileSize * tileScale,
    })),
  }),
  getAnimatedSpriteByName: function (
    names: string[]
  ): AnimatedSprite {
    return {
      sprites: names
        .map((name) => this.getTileByName(name))
        .filter((tile) => tile !== undefined),
    };
  },
});

export type Tiles = ReturnType<typeof tiles>;
