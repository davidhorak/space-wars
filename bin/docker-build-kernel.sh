docker build -t space-wars-kernel:latest --file=kernel/Dockerfile .
containerId=$(docker create space-wars-kernel:latest)
docker cp "$containerId":/kernel/space-wars.wasm ./client/public/space-wars.wasm 
docker rm "$containerId"