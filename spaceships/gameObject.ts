import {
  Asteroid,
  Explosion,
  GameObject,
  Laser,
  Rocket,
  Spaceship,
} from "./types";

export const isAsteroid = (gameObject: GameObject): gameObject is Asteroid =>
  gameObject.type === "asteroid";

export const isExplosion = (gameObject: GameObject): gameObject is Explosion =>
  gameObject.type === "explosion";

export const isLaser = (gameObject: GameObject): gameObject is Laser =>
  gameObject.type === "laser";

export const isRocket = (gameObject: GameObject): gameObject is Rocket =>
  gameObject.type === "rocket";

export const isSpaceship = (gameObject: GameObject): gameObject is Spaceship =>
  gameObject.type === "spaceship";
