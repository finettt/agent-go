FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -s" -o agent-go ./src

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN mkdir -p /home/appuser/.config/agent-go && chown -R appuser:appgroup /home/appuser/.config/agent-go
VOLUME /home/appuser/.config/agent-go

WORKDIR /app

COPY --from=builder /app/agent-go .
COPY AGENTS.md .
COPY .env /app/.env

RUN chown appuser:appgroup agent-go

USER appuser

ENTRYPOINT ["./agent-go"]