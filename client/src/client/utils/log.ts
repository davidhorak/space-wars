import type { CollisionLog, DamageLog, GameStateLog, KillLog, Log } from "../../../../spaceships";

export const isCollisionLog = (log: Log): log is CollisionLog => {
  return log.logType === "collision";
};

export const isGameStateLog = (log: Log): log is GameStateLog => {
  return log.logType === "game_state";
};

export const isKillLog = (log: Log): log is KillLog => {
  return log.logType === "kill";
};

export const isDamageLog = (log: Log): log is DamageLog => {
  return log.logType === "damage";
};
