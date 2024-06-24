PROJECT-NAME := l10n
GO := go
BUILD-DIR := build

ifeq ($(OS),Windows_NT)
	APP-PATH = .\cmd\$(PROJECT-NAME)
	OUTPUT-FILE = $(BUILD-DIR)\$(PROJECT-NAME).exe
	CLEAN = rmdir /s /q $(BUILD-DIR)
else
	APP-PATH = ./cmd/$(PROJECT-NAME)
	OUTPUT-FILE = $(BUILD-DIR)/$(PROJECT-NAME)
	CLEAN = rm -rf $(BUILD-DIR)
endif

.PHONY: build
build:
	$(GO) build -o $(OUTPUT-FILE) $(APP-PATH)

.PHONY: test
test:
	$(GO) test -v -coverpkg=./... ./...

.PHONY: clean
clean:
	$(CLEAN)