.DEFAULT_GOAL := build

# Detectar OS actual
ifeq ($(OS),Windows_NT)
CURRENT_OS := windows
else
CURRENT_OS := $(shell uname | tr '[:upper:]' '[:lower:]')
endif

# Detectar ARCH actual
ifeq ($(CURRENT_OS),windows)
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
CURRENT_ARCH := amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),ARM64)
CURRENT_ARCH := arm64
	endif
else
CURRENT_ARCH := $(shell uname -m)
# Convertir arquitectura a formato Go
	ifeq ($(CURRENT_ARCH),x86_64)
CURRENT_ARCH := amd64
	endif
	ifeq ($(CURRENT_ARCH),aarch64)
CURRENT_ARCH := arm64
	endif
endif

# Variables de versión
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Directorio de instalación
BINDIR ?= $(HOME)/bin

# Directorio de compilación
BUILD_DIR = build

# Nombre del binario
BINARY_NAME = rtg

# Módulo Go
GO_MODULE = github.com/didit-pub/rt-gcli

# Flags de compilación
LDFLAGS := -ldflags \
	"-X ${GO_MODULE}/pkg/version.Version=${VERSION} \
	-X ${GO_MODULE}/pkg/version.CommitSHA=${COMMIT_SHA} \
	-X ${GO_MODULE}/pkg/version.BuildDate=${BUILD_DATE}"

# Definir plataformas
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Directorio de distribución
DIST_DIR := dist

# Función para crear el directorio de distribución
.PHONY: $(DIST_DIR)
$(DIST_DIR):
	mkdir -p $@

# Función para obtener el nombre del binario
binary_name = $(DIST_DIR)/$(BINARY_NAME)-$(1)-$(2)$(if $(filter windows,$(1)),.exe,)

# Generar objetivos dinámicamente para cada plataforma
define PLATFORM_template
build-$(1)-$(2): $(DIST_DIR)
	@echo "Building for $(1)/$(2)..."
	@GOOS=$(1) GOARCH=$(2) go build $(LDFLAGS) -o $(call binary_name,$(1),$(2)) cmd/rtg/main.go
endef

# Generar objetivos para cada plataforma
$(foreach platform,$(PLATFORMS),\
	$(eval os := $(word 1,$(subst /, ,$(platform))))\
	$(eval arch := $(word 2,$(subst /, ,$(platform))))\
	$(eval $(call PLATFORM_template,$(os),$(arch))))

# Lista de todos los objetivos de compilación
BUILD_TARGETS := $(foreach platform,$(PLATFORMS),\
	build-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform))))

# Objetivo para limpiar el directorio de distribución
.PHONY: clean
clean:
	rm -rf $(DIST_DIR)

# Objetivo para compilar todos los binarios
.PHONY: build-all $(BUILD_TARGETS)
build-all: $(BUILD_TARGETS)

# Objetivo para mostrar plataformas disponibles
.PHONY: list-platforms
list-platforms:
	@echo "Plataformas disponibles:"
	@$(foreach platform,$(PLATFORMS),\
		echo "  $(platform)";)

# Objetivo para compilación individual
.PHONY: build
build: $(DIST_DIR) build-$(CURRENT_OS)-$(CURRENT_ARCH)
	mv $(DIST_DIR)/$(BINARY_NAME)-$(CURRENT_OS)-$(CURRENT_ARCH)$(if $(filter windows,$(CURRENT_OS)),.exe,) $(BINARY_NAME)

# test:
# 	go test ./...

# Objetivo para instalar el binario
.PHONY: install
install: build
	install -m 755 $(BINARY_NAME) $(BINDIR)/$(BINARY_NAME)

# Objetivo para crear una release
.PHONY: release
release:
	@if [ -z "$(NEW_VERSION)" ]; then \
		echo "Error: NEW_VERSION is not set. Usage: make release NEW_VERSION=X.Y.Z"; \
		exit 1; \
	fi
	git tag -a v$(NEW_VERSION) -m "Release version $(NEW_VERSION)"
	git push origin v$(NEW_VERSION)

.PHONY: pr-close
pr-close:
	@PR_ID=$$(gh pr view --json number -q '.number') && \
	gh pr merge -m -d $$PR_ID

.PHONY: pr
pr-%:
	git checkout -b $*
	git push --set-upstream origin $*
	gh pr create --fill
