import type {
    SpaceshipManager,
    SpaceshipManagerFactory,
    SpaceState,
    SpaceshipAction,
  } from "../..";
  
  const enhancedSpaceship = (): SpaceshipManager => {
    let engineThrustEnergyTrigger = 100;
    let attacking = true;
    let lastRocketFireMs = 0;
  
    const onUpdate = ({
      spaceship: self,
      deltaTimeMs,
      gameObjects,
    }: SpaceState): SpaceshipAction[] => {
      const actions: SpaceshipAction[] = [];
      if (!self) {
        return actions;
      }
  
      if (lastRocketFireMs > 0) {
        lastRocketFireMs -= deltaTimeMs;
        if (lastRocketFireMs <= 0) {
          lastRocketFireMs = 0;
        }
      }
  
      // Energy management: stop attacking if energy is below 30
      if (self.energy <= 30) {
        attacking = false; // Stop attacking
        engineThrustEnergyTrigger = 30 + Math.random() * 50; // Conserve more energy
      } else if (self.energy >= engineThrustEnergyTrigger) {
        engineThrustEnergyTrigger = Infinity;
        attacking = true;
        const mainThrust = 60 + Math.random() * 40; // Increase thrust efficiency
        let leftThrust = 0;
        let rightThrust = 0;
        if (Math.random() < 0.5) {
          leftThrust = 15 + Math.random() * 85; // Reduce randomness for more consistent movement
        } else {
          rightThrust = 15 + Math.random() * 85;
        }
        actions.push(["setEngineThrust", mainThrust, leftThrust, rightThrust]);
      }
  
      // Dodging logic using gameObjects data
      if (gameObjects) {
        const closestProjectile = gameObjects
          .filter((obj) => obj.type === "Projectile" || obj.type === "Asteroid")
          .reduce((closest, obj) => {
            const distance = Math.hypot(self.x - obj.x, self.y - obj.y);
            return distance < closest.distance
              ? { distance, obj }
              : closest;
          }, { distance: Infinity, obj: null });
  
        if (closestProjectile.obj?.type === "Asteroid" && closestProjectile.distance < 150) {
          // If the object is an Asteroid, the ship must take extra care to avoid it
          const dx = self.x - closestProjectile.obj.x;
          const dy = self.y - closestProjectile.obj.y;
          const evadeAngle = Math.atan2(dy, dx); // Calculate escape direction
          actions.push([
            "setEngineThrust",
            50,
            Math.cos(evadeAngle) * 90, // Move away from the asteroid
            Math.sin(evadeAngle) * 90,
          ]);
        } else if (closestProjectile.distance < 100) {
          let dodgeLeft = Math.random() > 0.5;
          actions.push([
            "setEngineThrust",
            50,
            dodgeLeft ? 90 : 0, // Sharp turn left or right for other objects
            dodgeLeft ? 0 : 90,
          ]);
        }
      }
  
      // Attacking logic
      if (attacking && self.energy > 30) { // Only attack if energy is above 30
        if (self.energy > 20 && self.rocketReloadTimerSec == 0 && self.rockets > 0) {
          actions.push(["fireRocket"]);
          lastRocketFireMs = 500;
        } else if (
          lastRocketFireMs <= 0 &&
          self.energy > 10 &&
          self.laserReloadTimerSec == 0
        ) {
          actions.push(["fireLaser"]);
        }
      }

      if(self.energy < 40){
        actions.push(["setEngineThrust", 0, 0, 0]);
      }
  
      return actions;
    };
  
    return {
      name: "TnT",
      onUpdate: onUpdate,
      onStart: function () {
        this.onReset();
      },
      onReset: () => {
        attacking = true;
        lastRocketFireMs = 0;
        engineThrustEnergyTrigger = 100;
      },
    };
  };
  
  export const createTnT = (): SpaceshipManagerFactory => () =>
    Promise.resolve(enhancedSpaceship());