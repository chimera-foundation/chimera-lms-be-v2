FROM golang:1.25-alpine AS builder

# Needed for installing Go modules
RUN apk add --no-cache git 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o chimera-api ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o chimera-seeder ./cmd/seed/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D chimerauser
WORKDIR /home/chimerauser

COPY --from=builder /app/chimera-api .
COPY --from=builder /app/chimera-seeder .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts ./scripts
RUN chmod +x ./scripts/entrypoint.sh

RUN chown -R chimerauser:chimerauser /home/chimerauser && \
    chmod +x ./chimera-api ./chimera-seeder

USER chimerauser
EXPOSE 8000

ENTRYPOINT ["./scripts/entrypoint.sh"]