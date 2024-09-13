## Tactable Space Wars
Players take on the role of elite engineers tasked with programming AI bots to pilot
futuristic spaceships in a high-stakes interstellar arena.

**The goal is simple**: code your bot to outmaneuver, outgun, and outsmart the competition,
becoming the last ship standing in a relentless battle for supremacy.

[ ] Add video and screenshots

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

In Space Wars, intelligence and innovation are your greatest weapons. Code wisely, pilot your ships with precision,
and claim your place among the stars!

> Updated: Sep 23, 2024

### Spaceship
- Starts with 100 points of health and 100 points of energy.
- Has 3 engines
 - Main truster (back)
 - 2 navigation trusters (sides)
 - Each could be operated independently
- Has rechargeable fuel. Using the engines depletes the energy. Not using the engines recharges the energy.
- Can shoot, laser beams - consumes energy.
- Can shoot rockets - consumes energy and has a limited amount.

#### Movement
- Each engine has its own throttle, from 0 to 100.
- The trust from all engines is added together to get the total trust.
- The trust is applied in the direction of the spaceship.
- Using the engines consumes energy, x/s for the main thruster, and y/s for the navigation thrusters.
- The maximum speed is **X** m/s.
- [ ] add the actual numbers
- [ ] add images

#### Laser
- Laser consumes **X** energy.
- Laser has reload time. **X**
- Laser has a damage radius. **X**
- Laser must hit precisely the target.
- Laser is faster than the rocket.
- [ ] add the actual numbers
- [ ] add images

#### Rockets
- Rockets consume **X** energy.
- Rockets have a limited amount. **10 per match**.
- Rockets have large damage trigger **X** radius, hence it can hit opponent in close proximity.
- Rockets have reload time. **X**
- Rockets cause **X** damage.
- [ ] add the actual numbers
- [ ] add images

#### Collisions
- Any object colliding with an asteroid is destroyed.
- A spaceship colliding with an opponent destroys both spaceships.
- A laser and a rocket launched do not collide with its launcher.

### Scoring
- Hitting an opponent with a rocket scores **X** points.
- Hitting an opponent with a laser scores **X** points.
- [ ] add the actual numbers


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

### How to Run the Game
- Go to [client](client) and run `yarn dev` to start the frontend.

---

### TODO
- [ ] Test the game and balance the game, mainly the truster physics
- [ ] Add all kernel tests
- [ ] Optimize the performance of the game
- [ ] Polish the game tiles
- [ ] [Spaceships TODOs](spaceships/readme.md)

---

## Contributing

See [contributing.md](contributing.md).

## License

Illogical is released under the MIT license. See [license.md](license.md).
