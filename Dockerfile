# ---- Frontend Build ----
FROM node:22-alpine AS frontend-build
WORKDIR /app/frontend
RUN npm install -g pnpm
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

# ---- Backend Build ----
FROM golang:1.26-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/main.go

# ---- Runtime ----
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=backend-build /app/server ./server
COPY --from=frontend-build /app/frontend/dist ./public

ENV APP_ENV=production
ENV HTTP_ADDR=:8080
ENV GIN_MODE=release

EXPOSE 8080
CMD ["./server"]
