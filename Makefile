.PHONY: all build compress run clean

BINARY_NAME=agent-go

all: build compress

build:
	@echo "Building..."
	@cd src && go build -ldflags="-s -w" -o ../$(BINARY_NAME) .

run: build
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)

compress:
	@echo "Compressing..."
	@upx --best --lzma ./$(BINARY_NAME)