FROM node:22-slim AS base
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

FROM base AS build-js
COPY . .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY --from=build-go /app/src/bin/app.wasm /app/src/bin/app.wasm
RUN pnpm run build

FROM gcr.io/distroless/nodejs22-debian12:nonroot
ENV NODE_ENV production

COPY --from=prod-deps /app/node_modules /app/node_modules
COPY --from=build-js /app/dist /app
COPY --from=build-go /app/src/bin /app/bin

EXPOSE 3000
CMD ["/app/index.js"]

