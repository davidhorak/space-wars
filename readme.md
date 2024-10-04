## Tactable Space Wars

Players take on the role of elite engineers tasked with programming AI bots to pilot
futuristic spaceships in a high-stakes interstellar arena.

**The goal is simple**: code your bot to outmaneuver, outgun, and outsmart the competition,
becoming the last ship standing in a relentless battle for supremacy.

https://github.com/user-attachments/assets/f78ffe1a-0e60-4d07-b52e-3ddbcc31a05e

### Game/Simulation Features:

- **Code Your Victory**: Players simply implement a [SpaceshipManager](spaceships/spaceshipManager.ts) interface to control their ship. The actual implementation can be written in any language that can compile to WASM, besides the pure TypeScript/JavaScript implementation.
- **Dynamic Arsenal**: Your ship is equipped with a variety of powerful weapons, including high-speed jet
  engines for evasive maneuvers, devastating laser beams for precision attacks, and rocket projectiles for
  explosive impact. Each weapon has its own strengths and weaknesses, requiring careful consideration in your
  bot's programming.
- **Hazardous Environment**: The battleground is strewn with asteroids, adding an extra layer of challenge.
  Bots must be programmed to navigate these obstacles skillfully to avoid catastrophic collisions that could end the match
  prematurely.
- **Real-Time Strategy**: The arena is a constantly evolving battlefield. As matches progress, players must adapt their bots'
  tactics to changing conditions and opponents' strategies, making quick thinking and adaptability key to victory.
- **Learning and Iteration**: After each match, players can analyze their bots' performance, tweak their code,
  and refine their strategies, fostering a cycle of learning and improvement.

![space-wars_screenshot](https://github.com/user-attachments/assets/802b0324-44bc-4659-8c83-2500029e43ba){: .mx-auto}

In Space Wars, intelligence and innovation are your greatest weapons. Code wisely, pilot your ships with precision,
and claim your place among the stars!

> Updated: Sep 23, 2024

> Notes:
>
> - 1 space unit is equal to **1** meter.
> - The actual numbers are not final and are subject to change. All configuration constants are located in [configuration.go](kernel/game/configuration.go)

### Spaceship

- Starts with **100 points** of health and **100 points** of energy.
- The size of the spaceship is **30** meters (radius).
- Has 3 engines
  - Main truster (back)
  - 2 navigation trusters (sides)
- Each could be operated independently
- Using the engines and firing lasers depletes the energy.
- Not using the engines recharges the energy. The recharge rate is **12.5** energy per second.
- Can shoot laser beams - consumes energy.
- Can shoot rockets - consumes energy and has a limited amount.

#### Movement

- Each engine has its own throttle, from **0** to **100**.
- The main thruster consumes **12.5** energy per second.
- The navigation thrusters consume **6.66** energy per second.
- The trust from all engines is added together to get the total trust.
- The side thrusters are **2** times weaker than the main thruster.
- The trust is applied in the direction of the spaceship.
- Using the engines consumes energy, x/s for the main thruster, and y/s for the navigation thrusters.
- The maximum speed is **192** m/s.
- The spaceship will reach the max speed in **5** seconds with full throttle.
- The spaceship will come close to a stop in **10** seconds without any thrust, caused by drag.
- The objects going over the screen wrap around to the other side.

https://github.com/user-attachments/assets/84e4e892-7ec1-4950-a1c2-72afd8f28de1

#### Laser

- Laser has a speed of **320** m/s.
- Laser consumes **6** energy.
- Laser has a lifespan of **5** seconds.
- Laser deals **20** damage.
- Laser must hit precisely the target.
- Laser has reload time. **250** milliseconds.

#### Rockets

- Rocket has a speed of **274** m/s.
- Rocket consume **20** energy.
- Rocket has a lifespan of **10** seconds.
- Rocket deals **60** damage to the target.
- Rocket has a **20** meters radius of explosion.
- Each ship has **10** rockets.
- Rocket has **1** second reload time.

### Collisions

- Any object colliding with an asteroid is destroyed.
- A spaceship colliding with an opponent destroys both spaceships.
- A laser and a rocket launched do not collide with its launcher.

### Start Location

- The spaceships are evenly distributed around the center of the map on a circle with a radius of 3/4 of the map width/height.
- The spaceships are randomly rotated.
- The position and rotation is randomized every match.

### Asteroids

- Asteroids are randomly generated on the map.
- The number of asteroids is randomized between **2** and **7**.
- The size of the asteroids is randomized between **10** and **30**.
- The minimum distance between asteroids is **10**.

### Scoring

- Hitting an opponent with a rocket scores **30** points.
- Hitting an opponent with a laser scores **10** points.
- Killing an opponent scores **100** points.

### Random Seed

- The seed is a 64-bit integer.
- The seed is used to generate the same map layout for the same seed.
- Seed could be provided via URL parameters `seed`, e.g. `localhost:3000/?seed=1234567890`

### Battlefield Size

- The width and height of the battlefield are set to **1024** (width) by **768** (height) meters.
- The width and height could be overridden via URL parameters `width` and `height`, e.g. `localhost:3000/?width=1200&height=800`

### FPS

- The FPS is set to **30**.
- The FPS could be overridden via URL parameters `fps`, e.g. `localhost:3000/?fps=60`

### Pause and Step

- The game could be paused and run step by step. Each step is one tick of the game with delta time of **50** milliseconds.

---

## How to Contribute (Spaceships)

Read [spaceships/readme.md](spaceships/readme.md) for more information on how to contribute to the spaceship codebase, adding new spaceships.

## Project Structure

- [kernel](kernel) - Game engine written in Go, compiled as WASM.
- [frontend](client) - React client for the game.
- [kernel client](client/src/client/) - Kernel client for the game, written in TypeScript, responsible for communication between the frontend and the kernel, and for rendering the game.
- [spaceships](spaceships) - Implementation of the SpaceshipManager interface for the game.
- [ui](ui) - UI source files for the game tiles.

### Local Setup

- [Install Go 1.22+](https://go.dev/doc/install)
- [Install Node 18+](https://nodejs.org/en/download)
- [Install Yarn](https://yarnpkg.com/getting-started/install)
- Go to [client](client) and run `yarn install` to install the dependencies.

### How to Build the Kernel WASM module

Linux/MacOS:

```sh
GOOS=js GOARCH=wasm go build -o ./client/public/space-wars.wasm
```

Windows:

```powershell
powershell -Command { $env:GOOS="js"; $env:GOARCH="wasm"; go build -o ./client/public/space-wars.wasm }
```

### How to run the kernel tests with full coverage

Linux/MacOS:

```sh
cd kernel
go test -coverprofile="coverage.out" ./... && go tool cover -html="coverage.out" -o coverage.html
```

Windows:

```powershell
cd kernel
powershell -Command { go test -coverprofile="coverage.out" ./...; go tool cover -html="coverage.out" -o coverage.html }
```

### How to Run the Game

- Go to [client](client) and run `yarn dev` to start the frontend.

---

### TODO

- [x] Add CI
- [ ] Test the game and balance the game, mainly the truster physics
- [x] Add all kernel tests
- [ ] Optimize the performance of the game
- [ ] Polish the game tiles
- [ ] Allow to load the game state form a json string
- [ ] Add a guide how to run the kernel without the UI client.
- [ ] [Spaceships TODOs](spaceships/readme.md)

---

## Contributing

See [contributing.md](contributing.md).

## License

Illogical is released under the MIT license. See [license.md](license.md).
