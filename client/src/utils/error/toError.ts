export const toError = (error: unknown): Error => {
  if (error instanceof Error) {
    return error;
  }
  return new Error(`${error}`);
};
