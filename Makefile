.PHONY: all build run clean

BINARY_NAME=agent-go

all: build

build:
	@echo "Building..."
	@cd src && go build -o ../$(BINARY_NAME) .

run: build
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)