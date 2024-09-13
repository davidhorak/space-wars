type SpriteKey = string
type SpriteImageUrl = string
type SpriteImageBase64Data = string
type SpriteImageSource = SpriteImageUrl | SpriteImageBase64Data
type FromBlob = boolean

export const sprites = () =>
  ((cache) => ({
    load: (...sprites: [SpriteKey, SpriteImageSource, FromBlob][]) =>
      new Promise<void>((resolve) => {
        let pending = sprites.length
        for (const [key, data, fromBlob] of sprites) {
          const img = new window.Image()
          img.onload = (): void => {
            cache.set(key, img)
            pending--
            if (pending === 0) {
              resolve()
            }
          }

          img.src = fromBlob ? `data:image/png;base64,${data}` : data
        }
      }),
    get: (key: string) => {
      const sprite = cache.get(key)
      if (!sprite) {
        throw new Error(`missing sprint ${key}`)
      }
      return sprite
    }
  }))(new Map<string, HTMLImageElement>())