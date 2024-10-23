// Relative path to go-wasm-exec.js file
import "../../client/public/go-wasm-exec.js";
import fs from "fs";

async function main() {
  const go = new Go();
  const result = await WebAssembly.instantiate(
    // Relative path to kernel wasm file (see https://github.com/davidhorak/space-wars?tab=readme-ov-file#how-to-build-the-kernel-wasm-module)
    fs.readFileSync("../../client/public/space-wars.wasm"),
    go.importObject
  );
  go.run(result.instance);

  // The kernel interface lives in the global `spaceWars` object
  // To see the full documentation of the kernel interface, see:
  // https://github.com/davidhorak/space-wars/blob/master/client/src/vite-env.d.ts

  // Battlefield size
  const width = 1024;
  const height = 768;

  // 1. Initialize the kernel
  spaceWars.init(width, height);
  // Optionally pass a seed to the init function
  // spaceWars.init(width, height, seed);

  // 2. Add spaceships
  // There must be at least two spaceships
  spaceWars.addSpaceship(
    "Spaceship 1",
    100, // X
    100, // Y
    3.14 // Angle in radians
  );
  spaceWars.addSpaceship(
    "Spaceship 2",
    200, // X
    200, // Y
    3.14 // Angle in radians
  );

  // 3. Start the game
  spaceWars.start();
  let gameState = spaceWars.state();
  console.log(gameState.status);

  // 4.1 Run the game loop
  let iteration = 0;
  while (iteration < 10) {
    spaceWars.tick(
      200 // delta time in milliseconds
    );

    // 4.2 Get the game state
    // The full documentation of the game state can be found:
    // https://github.com/davidhorak/space-wars/blob/master/spaceships/types.ts#L154
    gameState = spaceWars.state();

    // Use the game state to control the spaceships
    // For example, we can check the position of the spaceships
    const spaceship = gameState.gameObjects.find(
      (obj) => obj.name === "Spaceship 1"
    );
    console.log(spaceship.position);

    // 4.3 Perform actions
    // The full documentation of the actions can be found:
    // https://github.com/davidhorak/space-wars/blob/master/spaceships/spaceshipAction.ts
    spaceWars.action(
      "setEngineThrust", // Action type
      "Spaceship 1", // Spaceship name
      100, // Main engine thrust
      0, // Left engine thrust
      0 // Right engine thrust
    );
  
    iteration++;
  }
}

await main();
