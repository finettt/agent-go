FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -s" -o agent-go ./src

FROM alpine:latest

RUN apk add nodejs
RUN addgroup -S agent-go && adduser -S finett -G agent-go

RUN mkdir -p /home/finett/.config/agent-go && chown -R finett:agent-go /home/finett/.config/agent-go
VOLUME /home/finett/.config/agent-go

# Create workspace directory for mounting host files
RUN mkdir -p /workspace && chown -R finett:agent-go /workspace
VOLUME /workspace

WORKDIR /app

COPY --from=builder /app/agent-go .

RUN chown finett:agent-go agent-go

USER finett

# Set working directory to workspace so mounted files are accessible
WORKDIR /workspace

ENTRYPOINT ["/app/agent-go"]