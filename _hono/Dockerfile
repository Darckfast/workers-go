FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM golang:1.21 AS build-go
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=js GOARCH=wasm go build -ldflags='-s -w' -o ./src/bin/app.wasm main.go

FROM base AS build
COPY . .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY --from=build-go /app/src/bin/app.wasm /app/src/bin/app.wasm
RUN pnpm run build

FROM base
RUN apt update && apt install curl -y

# COPY migrations migrations

COPY package.json package.json
COPY --from=prod-deps /app/node_modules node_modules
COPY --from=build /app/dist ./dist
COPY --from=build-go /app/src/bin ./dist/bin

EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=30s --retries=3 --start-period=3s --start-interval=150ms CMD curl -f http://localhost:3000/health || exit 1
CMD ["pnpm", "start"]

