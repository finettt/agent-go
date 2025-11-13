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

# Create workspace directory for mounting host files
RUN mkdir -p /workspace && chown -R appuser:appgroup /workspace
VOLUME /workspace

WORKDIR /app

COPY --from=builder /app/agent-go .

RUN chown appuser:appgroup agent-go

USER appuser

# Set working directory to workspace so mounted files are accessible
WORKDIR /workspace

ENTRYPOINT ["./agent-go"]