FROM node:22-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app

FROM base AS prod-deps
COPY . .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod

FROM golang:1.23 AS build-go
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=js GOARCH=wasm go build -o ./bin/app.wasm ./pkg/main.go

FROM base AS build
COPY . .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install
COPY --from=build-go /app/bin/app.wasm /app/bin/app.wasm
RUN pnpm run build:js

FROM base
COPY migrations migrations
COPY package.json package.json
COPY wrangler.toml wrangler.toml
COPY --from=prod-deps /app/node_modules node_modules
COPY --from=build /app/dist .
COPY --from=build-go /app/bin bin

EXPOSE 5173
RUN pnpm add wrangler
CMD ["pnpm", "start"]
