FROM node:22-slim AS base
RUN apt update && apt install git -y

FROM base AS deps

RUN corepack enable
WORKDIR /app
# COPY pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch
COPY package.json ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install --prod

FROM golang:1.21 AS build-go
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=js GOARCH=wasm go build -ldflags='-s -w' -o ./bin/app.wasm main.go

FROM base AS build
RUN corepack enable

WORKDIR /app
# COPY pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch
COPY package.json ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install
COPY . .
RUN pnpm build

FROM base

WORKDIR /app
COPY --from=deps /app/node_modules /app/node_modules
COPY --from=build /app/dist /app/dist
COPY --from=build-go /app/bin /app/dist/bin
ENV NODE_ENV production
CMD ["node", "./dist/index.js"]
