export type SquareCollider = {
  type: "square";
  enabled: boolean;
  position: {
    x: number;
    y: number;
  };
  rotation: number;
  size: {
    width: number;
    height: number;
  };
};

export type CircleCollider = {
  type: "circle";
  enabled: boolean;
  position: {
    x: number;
    y: number;
  };
  radius: number;
};

export type PolygonCollider = {
  type: "polygon";
  enabled: boolean;
  position: {
    x: number;
    y: number;
  };
  rotation: number;
  vertices: { x: number; y: number }[];
};

export type GameObject = {
  id: number;
  type: string;
  enabled: boolean;
  position: {
    x: number;
    y: number;
  };
};

export type Asteroid = GameObject & {
  type: "asteroid";
  enabled: boolean;
  radius: number;
  lifespanSec: number;
  collider: CircleCollider;
};

export type Explosion = GameObject & {
  type: "explosion";
  enabled: boolean;
  radius: number;
  durationSec: number;
  lifespanSec: number;
};

export type Projectile = GameObject & {
  enabled: boolean;
  rotation: number;
  velocity: {
    x: number;
    y: number;
  };
  lifespanSec: number;
  damage: number;
  owner: string;
};

export type Laser = Projectile & {
  type: "laser";
  collider: SquareCollider;
};

export type Rocket = Projectile & {
  type: "rocket";
  collider: CircleCollider;
};

export type Spaceship = GameObject & {
  type: "spaceship";
  enabled: boolean;
  destroyed: boolean;
  name: string;
  startPosition: {
    x: number;
    y: number;
  };
  rotation: number;
  velocity: {
    x: number;
    y: number;
  };
  health: number;
  energy: number;
  engine: {
    mainThrust: number;
    leftThrust: number;
    rightThrust: number;
  };
  rockets: number;
  kills: number;
  score: number;
  laserReloadTimerSec: number;
  rocketReloadTimerSec: number;
  collider: CircleCollider;
};

export type Log = {
  id: number;
  logType: "damage" | "kill" | "collision" | "game_state";
  message: string;
  time: string;
  meta: Record<string, string>;
};

export type DamageLog = Log & {
  logType: "damage";
  meta: {
    who: string;
    whom: string;
    damage: string;
    damageType: string;
  };
};

export type KillLog = Log & {
  logType: "kill";
  meta: {
    who: string;
    whom: string;
  };
};

export type CollisionLog = Log & {
  logType: "collision";
  meta: {
    who: string;
    with: string;
  };
};

export type GameStateLog = Log & {
  logType: "game_state";
  meta: {
    status: "initialized" | "running" | "paused" | "ended";
  };
};

export type GameState = {
  status: "initialized" | "running" | "paused" | "ended";
  seed: number;
  size: {
    width: number;
    height: number;
  };
  gameObjects: (Asteroid | Explosion | Laser | Rocket | Spaceship)[];
  logs: Log[];
};
