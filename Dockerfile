FROM node:22-alpine AS web-builder

WORKDIR /src/web/app

COPY web/app/package*.json ./
RUN npm ci

COPY web/app/ ./
RUN npm run build


FROM golang:1.24-alpine AS go-builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mtesense-home ./cmd/server


FROM alpine:3.22

WORKDIR /app

ENV PORT=8080
ENV DATABASE_PATH=/app/data/app.db
ENV UPLOAD_DIR=/app/public_uploads

RUN mkdir -p /app/data /app/public_uploads /app/web/app/dist

COPY --from=go-builder /out/mtesense-home /app/mtesense-home
COPY --from=web-builder /src/web/app/dist /app/web/app/dist

EXPOSE 8080

CMD ["/app/mtesense-home"]
