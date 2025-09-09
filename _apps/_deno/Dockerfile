FROM golang:1.21 AS build-go
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=js GOARCH=wasm go build -ldflags='-s -w' -o app.wasm main.go

FROM denoland/deno:2.4.5

# Create working directory
WORKDIR /app

# Copy source
COPY . .

# Compile the main app
RUN deno cache main.ts

# Copy the compiled wasm binary
COPY --from=build-go /app/app.wasm ./bin/app.wasm

# Run the app
CMD ["deno", "run", "--allow-net", "--allow-read", "main.ts"]

