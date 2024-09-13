export const observable = <T>() =>
  ((observers) => ({
    subscribe: (observer: (subject: T) => void) => {
      const key = Symbol();
      observers.set(key, observer);
      return key;
    },
    unsubscribe: (key: symbol) => observers.delete(key),
    broadcast: (subject: T): void => {
      for (const observer of observers.values()) {
        observer(subject);
      }
    },
  }))(new Map<symbol, (subject: T) => void>());
