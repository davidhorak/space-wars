import {
  SpaceshipManager,
  SpaceshipManagerFactory,
  SpaceState,
  SpaceshipAction,
} from "../..";

const lasterTester = (name: string): SpaceshipManager => {
  const onUpdate = ({ spaceship: self }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (self.energy > 10 && self.laserReloadTimerSec == 0) {
      actions.push(["fireLaser"]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

const rocketTester = (name: string): SpaceshipManager => {
  const onUpdate = ({ spaceship: self }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (self.energy > 20 && self.rocketReloadTimerSec == 0) {
      actions.push(["fireRocket"]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

const engineTester = (
  name: string,
  mainEngineThrust: number,
  leftEngineThrust: number,
  rightEngineThrust: number
): SpaceshipManager => {
  const onUpdate = ({ spaceship: self }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    if (self.energy >= 100) {
      actions.push([
        "setEngineThrust",
        mainEngineThrust,
        leftEngineThrust,
        rightEngineThrust,
      ]);
    }

    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

const rocinante = (name: string): SpaceshipManager => {

  let mode = "Attack";
  let laserBurst = 2;
  let target = -1;

  const onUpdate = ({ spaceship: self, deltaTimeMs, gameObjects }: SpaceState): SpaceshipAction[] => {
    const actions: SpaceshipAction[] = [];
    if (!self) {
      return actions;
    }

    //Target Lock
    if (target == -1) {
      gameObjects.forEach((gameObject, index) => {
        if (gameObject.type == "spaceship" && gameObject.id != self.id && !gameObject.destroyed){
          target = index;
        }
      });
    } else {
      let targetShip = gameObjects[target];
      if (targetShip.type == "spaceship") {
        if(targetShip.destroyed) {
          gameObjects.forEach((gameObject, index) => {
            if (gameObject.type == "spaceship" && gameObject.id != self.id && !gameObject.destroyed){
              target = index;
            }
          });    
        }        
      }
    }

    // Calculate the differences in position
    const dy = gameObjects[target].position.y - self.position.y;
    const dx = gameObjects[target].position.x - self.position.x;
  
    // Calculate the angle to obj2 from obj1
    const theta2 = Math.atan2(dy, dx); // Angle in radians

    // Calculate the direction difference
    const directionDifference = theta2 - self.rotation;
  
    // Normalize the angle difference to the range of -π to π
    const normalizedDifference = ((directionDifference + Math.PI) % (2 * Math.PI)) - Math.PI;

    if (self.energy == 0) {
      mode = "Recharge"
      let mainThrust = 0;
      let leftThrust = 0;
      let rightThrust = 0;
      actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]); 
    } else if (self.energy == 100) {
      mode = "Attack"
    }

    if (mode == "Attack") {
      if (normalizedDifference < 0.01 && normalizedDifference > -0.01) {
        actions.push(["fireLaser"]);
        let mainThrust = 0;
        let leftThrust = 0;
        let rightThrust = 0;
        actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);  
      } else {
        if (normalizedDifference > 0) {
          let mainThrust = 100;
          let leftThrust = 100;
          let rightThrust = 0;
          actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);  
        } else if (normalizedDifference < 0) {
          let mainThrust = 100;
          let leftThrust = 0;
          let rightThrust = 100;
          actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);  
    
        }
    }

 
    }

    // TODO:
    // Switch targets
    // Avoid asteroids
    // Leading shots
 
    // if (self.energy >= 70){
    //   mode = "Spin"
    //   laserBurst = 2;
    //   ;

    // } else {
    //   mode = "Off"
    // }

    
    // if (mode == "Spin") {
    //   let mainThrust = 100;
    //   let leftThrust = 100;
    //   let rightThrust = 0;
    //   actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);
    // } else if (mode == "Off") {
    //   let mainThrust = 0;
    //   let leftThrust = 0;
    //   let rightThrust = 0;
    //   actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);
    // }

    // if (laserBurst > 0 && self.laserReloadTimerSec == 0) {
    //   laserBurst = laserBurst - 1;
    //   actions.push(["fireLaser"]);
    // }
  
  


    return actions;
  };

  return {
    name,
    onUpdate: onUpdate,
    onStart: () => {},
    onReset: () => {},
  };
};

export const createRocinante =
  (name: string): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(rocinante(name));

export const createLaserTester =
  (name: string): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(lasterTester(name));

export const createRocketTester =
  (name: string): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(rocketTester(name));

export const createEngineTester =
  (
    name: string,
    mainEngineThrust: number,
    leftEngineThrust: number,
    rightEngineThrust: number
  ): SpaceshipManagerFactory =>
  () =>
    Promise.resolve(
      engineTester(name, mainEngineThrust, leftEngineThrust, rightEngineThrust)
    );

  
