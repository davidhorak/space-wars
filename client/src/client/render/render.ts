import { Vector2 } from "../utils";
import { Sprite } from "./Sprite";

const measureText = ((cache) => {
  return (text: string, font: string, context: CanvasRenderingContext2D): number => {
    const key = `${text}|${font}`;
    if (cache.has(key)) {
      return cache.get(key) ?? 0;
    }
    const width = context.measureText(text).width;
    cache.set(key, width);
    return width;
  };
})(new Map<string, number>())

export const render = (
  canvas: HTMLCanvasElement,
  context: CanvasRenderingContext2D
) => ({
  clear: () => context.clearRect(0, 0, canvas.width, canvas.height),

  createPattern: (sprite: Sprite): CanvasPattern => {
    return context.createPattern(sprite.image, 'repeat')!
  },

  drawPattern: (pattern: CanvasPattern, x: number, y: number, width: number, height: number) => {
    context.fillStyle = pattern
    context.fillRect(x, y, width, height)
  },

  drawSprite: (sprite: Sprite, x: number, y: number, width: number, height: number, rotation: number = 0) => {
    if (!sprite) {
      return;
    }
    
    let dx = x;
    let dy = y;

    if (rotation !== 0) {
        context.translate(x, y);
        context.rotate(rotation);
        dx = -width / 2;
        dy = -height / 2;
    }

    context.drawImage(
      sprite.image,
      sprite.x,
      sprite.y,
      sprite.width,
      sprite.height,
      dx,
      dy,
      width,
      height
    )

    if (rotation !== 0) {
        context.rotate(-rotation);
        context.translate(-x, -y);
    }
  },

  drawRect: (
    color: string,
    lineWidth: number,
    x: number,
    y: number,
    width: number,
    height: number,
    rotation: number = 0
  ) => {
    let dx = x;
    let dy = y;

    if (rotation !== 0) {
        context.translate(x, y);
        context.rotate(rotation);
        dx = -width / 2;
        dy = -height / 2;
    }
  
    context.beginPath()
    context.lineWidth = lineWidth
    context.strokeStyle = color
    context.rect(dx, dy, width, height)
    context.stroke()

    if (rotation !== 0) {
        context.rotate(-rotation);
        context.translate(-x, -y);
    }
  },

  drawRectFilled: (
    color: string | CanvasGradient | CanvasPattern,
    x: number,
    y: number,
    width: number,
    height: number
  ) => {
    context.fillStyle = color;
    context.fillRect(x, y, width, height);
  },

  drawCircle: (color: string, lineWidth: number, x: number, y: number, radius: number) => {
    context.beginPath()
    context.lineWidth = lineWidth
    context.strokeStyle = color
    context.arc(x, y, radius, 0, 2 * Math.PI)
    context.stroke()
  },

  drawCircleFilled: (
    color: string | CanvasGradient | CanvasPattern,
    x: number,
    y: number,
    radius: number
  ) => {
    context.beginPath()
    context.arc(x, y, radius, 0, 2 * Math.PI)
    context.fillStyle = color
    context.fill()
  },

  drawPolygon: (color: string, lineWidth: number, points: Vector2[]) => {
    context.lineWidth = lineWidth
    context.strokeStyle = color

    context.beginPath();
    const [start, ...rest] = points

    context.moveTo(start.x, start.y);
    for (const point of rest) {
      context.lineTo(point.x, point.y);
    }

    context.closePath();
    context.stroke();
  },

  drawText: (font: string, color: string, text: string, x: number, y: number, centered = false) => {
    context.font = font
    context.fillStyle = color
    if (centered) {
      x -= measureText(text, font, context) / 2
    }
    context.fillText(text, x, y)
  },

  translate: (x: number, y: number) => {
    context.translate(x, y);
  },

  rotate: (radians: number) => {
    context.rotate(radians);
  },
});

export type Render = ReturnType<typeof render>;