## Run the kernel without the UI
A quick way to run the kernel without the UI is to use a node script. This is useful if you want to run the kernel in a node environment without the UI.

1. Build the kernel

See [How to build the kernel WASM module](https://github.com/davidhorak/space-wars?tab=readme-ov-file#how-to-build-the-kernel-wasm-module)

2. Create a node module script to run the kernel

See the [index.mjs](./index.mjs) file in this directory.

3. Run the script

```bash
node index.mjs
```
