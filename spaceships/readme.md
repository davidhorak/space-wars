## Space Wars Spaceships Manager
This is the space where you can implement your own spaceship manager. It's a crucial part of the Space Wars game, allowing you to control and manage your spaceship's actions and behaviors.

### Where to Store your Code?
Create a folder in the `/spaceships` directory with your GitHub username, for example: `/spaceships/john-doe/`. This will help keep your code organized and easily identifiable.

### How to Merge your Code?
Create a pull request to merge your code into the main branch. This will allow us to review your code and merge it into the game. 

### Development
- Implement the [SpaceshipManagerFactory](./spaceshipManager.ts) interface in your code. This factory function is responsible for creating instances of your spaceship manager.

```ts
export type SpaceshipManagerFactory = () => Promise<SpaceshipManager>;
```

- Add your spaceship manager to the [spaceships](./index.ts) array. This is where all spaceship managers are registered and made available for use in the game.

### SpaceshipManager Interface
```ts
export interface SpaceshipManager {
  name: string;
  onUpdate(state: SpaceState): SpaceshipAction[];
  onStart(width: number, height: number): void;
  onReset(): void;
}
```

#### Name
The name of your spaceship, it must be unique. This is used to identify your spaceship in the game.

#### onUpdate
This function is called every game tick, contains the current game state and expects the actions to be performed by the spaceship.

- **deltaTimeMs**: The time in milliseconds since the last update. This can be used to calculate the spaceship's movement and other time-dependent actions.
- **spaceship**: The spaceship object, contains the spaceship's current state. This includes its position, velocity, and other relevant properties.
- **gameObjects**: The game objects, contains all the game objects in the game. This includes asteroids, other spaceships, and any other objects that can interact with your spaceship.

#### onStart
This function is called when the game starts, it can be used to initialize any resources needed by the spaceship. This is a good place to set up any initial conditions or configurations for your spaceship. 

- **width**: The width of the game area.
- **height**: The height of the game area.

#### onReset
This function is called when the game resets, it can be used to reset any resources needed by the spaceship. This is useful for restarting the game or resetting the spaceship to its initial state.

## Spaceship Actions
### Set Engine Thrust
Example: 
```ts
["setEngineThrust", MainEngineThrust, LeftEngineThrust, RightEngineThrust]
```

- **MainEngineThrust**: The main engine thrust, it must be a value between 0 and 100. This controls the forward movement of the spaceship.
- **LeftEngineThrust**: The left engine thrust, it must be a value between 0 and 100. This controls the left movement of the spaceship.
- **RightEngineThrust**: The right engine thrust, it must be a value between 0 and 100. This controls the right movement of the spaceship.


### Fire Laser
Example:
```ts
["fireLaser"]
```

Fires a laser from the spaceship. This action can be used to attack other spaceships or asteroids.

### Fire Rocket
Example:
```ts
["fireRocket"]
```

Fires a rocket from the spaceship. This action can be used to attack other spaceships or asteroids with a more powerful projectile.

---

## Notes
`_samples` folder contains the examples and tests of the spaceships that you can use as a template to create your own spaceship manager.

## TODO
- [ ] Add more examples
- [ ] Add template for Python WASM, and any other languages that you want to support.
