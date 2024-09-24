export const Path = {
  Battlefield: '/',
  PageNotFound: '*',
} as const;

export type PathKeys = keyof typeof Path;
export type PathValues = typeof Path[PathKeys];
