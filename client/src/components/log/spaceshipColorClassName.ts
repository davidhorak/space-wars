const spaceshipColorClassNameMap = new Map<string, string>();

export const spaceshipColorClassName = (name: string, colors: number) => {
  if (spaceshipColorClassNameMap.has(name)) {
    return spaceshipColorClassNameMap.get(name);
  }

  const index = spaceshipColorClassNameMap.size % colors;
  const className = `color-palette-${index}`;
  spaceshipColorClassNameMap.set(name, className);
  return className;
};