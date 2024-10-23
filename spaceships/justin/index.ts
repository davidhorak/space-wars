import type {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
  GameState,
  Spaceship,
} from "../";
import { Vector2 } from "../../client/src/client/utils";
import { SetEngineThrustAction } from "../spaceshipAction";

const LOOK_AHEAD_DISTANCE = 100;

class DangerZone {
  public left: Vector2 = { x: 0, y: 0 };
  public right: Vector2 = { x: 0, y: 0 };
  public middle: Vector2 = { x: 0, y: 0 };

  constructor(spaceship: Spaceship) {
    this.middle = {
      x:
        spaceship.position.x +
        LOOK_AHEAD_DISTANCE * Math.cos(spaceship.rotation),
      y:
        spaceship.position.y +
        LOOK_AHEAD_DISTANCE * Math.sin(spaceship.rotation),
    };
    this.left = {
      x:
        this.middle.x +
        LOOK_AHEAD_DISTANCE * Math.cos(spaceship.rotation - Math.PI / 5),
      y:
        this.middle.y +
        LOOK_AHEAD_DISTANCE * Math.sin(spaceship.rotation - Math.PI / 5),
    };
    this.right = {
      x:
        this.middle.x +
        LOOK_AHEAD_DISTANCE * Math.cos(spaceship.rotation + Math.PI / 5),
      y:
        this.middle.y +
        LOOK_AHEAD_DISTANCE * Math.sin(spaceship.rotation + Math.PI / 5),
    };
  }

  // returns true if the object is in the danger zone
  // and if the object is on the left or right side of the danger zone
  public isObjectInDangerZone(
    spaceship: Spaceship,
    object: GameState["gameObjects"][number],
    dangerZone: DangerZone
  ): { isInDanger: boolean; side: "left" | "right" | null } {
    const { left, right, middle } = dangerZone;
    const { position } = spaceship;

    const objectPos = object.position;

    const crossProduct = (v1: Vector2, v2: Vector2) =>
      v1.x * v2.y - v1.y * v2.x;

    const vectorToLeft = { x: left.x - position.x, y: left.y - position.y };
    const vectorToRight = { x: right.x - position.x, y: right.y - position.y };
    const vectorToObject = {
      x: objectPos.x - position.x,
      y: objectPos.y - position.y,
    };

    const leftCross = crossProduct(vectorToLeft, vectorToObject);
    const rightCross = crossProduct(vectorToObject, vectorToRight);

    const isInDanger = leftCross > 0 || rightCross > 0;

    let side: "left" | "right" | null = null;
    if (isInDanger) {
      const vectorToMiddle = {
        x: middle.x - position.x,
        y: middle.y - position.y,
      };
      const middleCross = crossProduct(vectorToMiddle, vectorToObject);
      side = middleCross > 0 ? "right" : "left";
    }

    return { isInDanger, side };
  }
}

const justin = (name: string): SpaceshipManager => {
  let engineThrustEnergyTrigger = 100;
  let attacking = true;
  let lastRocketFireMs = 0;
  let lineOfSight: Vector2 = { x: 0, y: 0 };
  // let dangerZone;

  const getDistance = (position1: Vector2, position2: Vector2) => {
    return Math.sqrt(
      (position1.x - position2.x) ** 2 + (position1.y - position2.y) ** 2
    );
  };

  const handleFlightPath = (
    spaceship: Spaceship,
    closestFlightThreat: GameState["gameObjects"][number] | null,
    dangerZone: DangerZone
  ): SetEngineThrustAction => {
    let mainThrust = 0;
    let leftThrust = 0;
    let rightThrust = 0;
    if (spaceship.energy > 70) {
      if (closestFlightThreat === null) {
        if (spaceship.energy > 70) {
          mainThrust = 25 + Math.random() * 50;
        } else {
          mainThrust = 0;
        }
        leftThrust = 0;
        rightThrust = 0;
      } else {
        const { side } = dangerZone.isObjectInDangerZone(
          spaceship,
          closestFlightThreat,
          dangerZone
        );
        mainThrust = 25 + Math.random() * 50;
        if (side === "left") {
          leftThrust = 25 + Math.random() * 50;
        } else if (side === "right") {
          rightThrust = 25 + Math.random() * 50;
        } else {
          mainThrust = 0;
          leftThrust = 0;
          rightThrust = 0;
        }
      }
    }
    return ["setEngineThrust", mainThrust, leftThrust, rightThrust];
  };

  const onUpdate = ({
    spaceship: self,
    deltaTimeMs,
    gameObjects,
  }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    lineOfSight = {
      x: self.position.x + Math.cos(self.rotation),
      y: self.position.y + Math.sin(self.rotation),
    };

    const dangerZone = new DangerZone(self);

    const flightThreats = gameObjects.filter(
      (object) =>
        dangerZone.isObjectInDangerZone(self, object, dangerZone).isInDanger
    );

    let closestFlightThreat: GameState["gameObjects"][number] | null = null;
    flightThreats.forEach((flightThreat) => {
      if (!closestFlightThreat) {
        closestFlightThreat = flightThreat;
      }
      const distance = getDistance(self.position, flightThreat.position);
      if (distance < getDistance(self.position, closestFlightThreat.position)) {
        closestFlightThreat = flightThreat;
      }
    });

    // console.log("closestFlightThreat", closestFlightThreat);

    // if (flightThreats.length > 0) {
    //   console.log(flightThreats);
    // }

    // console.log("gameObjects", gameObjects);

    if (lastRocketFireMs > 0) {
      lastRocketFireMs -= deltaTimeMs;
      if (lastRocketFireMs <= 0) {
        lastRocketFireMs = 0;
      }
    }

    if (self.energy <= 10) {
      attacking = false;
      engineThrustEnergyTrigger = 25 + Math.random() * 50;
    }

    // if (self.energy >= 30) {
    //   engineThrustEnergyTrigger = Infinity;
    const flightPathAction = handleFlightPath(
      self,
      closestFlightThreat,
      dangerZone
    );
    //   console.log("flightPathAction", flightPathAction);
    actions.push(flightPathAction);
    // }
    if (self.energy >= 30) {
      attacking = true;
    }
    if (
      attacking &&
      self.energy > 60 &&
      self.rocketReloadTimerSec == 0 &&
      self.rockets > 0
    ) {
      // actions.push(["fireRocket"]);
      lastRocketFireMs = 500;
    } else if (
      lastRocketFireMs <= 0 &&
      attacking &&
      self.energy > 10 &&
      self.laserReloadTimerSec == 0
    ) {
      // actions.push(["fireLaser"]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: function (spaceship, width, height) {
      this.onReset(spaceship, width, height);
    },
    onReset: () => {
      attacking = true;
      lastRocketFireMs = 0;
      engineThrustEnergyTrigger = 100;
    },
  };
};

const name = "JUSTIN";

export const createJustin = (): SpaceshipManagerFactory => () =>
  Promise.resolve(justin(name));
