export interface Sprite {
  image: HTMLOrSVGImageElement;
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface AnimatedSprite {
  sprites: Sprite[];
}
