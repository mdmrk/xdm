GO := go build
GOFLAGS := -ldflags "-s -w"

BIN_PATH := bin
SRC_PATH := .

TARGET_SERVER_NAME := xdmediad
LOGGER_SERVER_NAME := xdmediad

TARGET_SERVER := $(BIN_PATH)/$(TARGET_SERVER_NAME)
LOGGER_SERVER := $(BIN_PATH)/$(LOGGER_SERVER_NAME)

OBJ_SERVER := $(SRC_PATH)/server/*.go
OBJ_LOGGER := $(SRC_PATH)/logger/*.go

CLEAN_LIST := $(BIN_PATH)

default: makedir all

.PHONY: makedir
makedir:
	@mkdir -p $(BIN_PATH)

.PHONY: all
all: server logger

.PHONY: server
server: $(OBJ_SERVER)
	$(GO) $(GOFLAGS) -o $(TARGET_SERVER) $(OBJ_SERVER)

.PHONY: server_run
server_run: server
	@$(TARGET_SERVER)

.PHONY: logger
logger: $(OBJ_LOGGER)
	$(GO) $(GOFLAGS) -o $(LOGGER_SERVER) $(OBJ_LOGGER)

.PHONY: logger_run
logger_run: logger
	@$(LOGGER_SERVER)

.PHONY: clean
clean:
	@echo CLEAN $(CLEAN_LIST)
	@rm -rf $(CLEAN_LIST)
