FROM node:22-slim AS base
RUN apt update && apt install git -y
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app
COPY . .

# Build the wasm binary
FROM golang:1.21 AS build-go
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=js GOARCH=wasm go build -o ./bin/app.wasm main.go

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM gcr.io/distroless/nodejs22-debian12:nonroot
ENV NODE_ENV=production

WORKDIR /app
COPY --from=prod-deps /app/node_modules ./node_modules
COPY --from=build /app/dist ./
COPY --from=build-go /app/bin ./bin
EXPOSE 5173
CMD ["index.js"]
