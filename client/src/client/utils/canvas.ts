const isHTMLCanvasElement = (
  element: HTMLElement | null
): element is HTMLCanvasElement => element?.tagName.toLowerCase() === "canvas";

export const getCanvas = (canvasId: string): HTMLCanvasElement => {
  const canvas = document.getElementById(canvasId);
  if (!isHTMLCanvasElement(canvas)) {
    throw new Error(`Game canvas #${canvasId} not found`);
  }
  return canvas;
};
