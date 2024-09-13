import { observable } from "./observable";

export const loop = ({ fps }: { fps: number }) => {
  const onUpdate = observable<number>()
  const onRender = observable()
  const fpMs = 1000 / fps
  let lastUpdateTimeMs = performance.now()
  let renderTimerMs = fpMs

  const loop = (timeMs: number) => {
    const deltaTimeMs = timeMs - lastUpdateTimeMs
    lastUpdateTimeMs = timeMs
    renderTimerMs -= deltaTimeMs
    onUpdate.broadcast(deltaTimeMs)

    if (renderTimerMs <= 0) {
      onRender.broadcast(deltaTimeMs)
      renderTimerMs = fpMs
    }

    window.requestAnimationFrame(loop)
  }

  window.requestAnimationFrame(loop)

  return {
    onUpdate: (callback: (deltaTimeMs: number) => void) => onUpdate.subscribe(callback),
    onRender: (callback: () => void) => onRender.subscribe(callback),
  }
};

export type Loop = ReturnType<typeof loop>;
