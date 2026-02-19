# ─────────────────────────────────────────────────────────────────────────────
#  Makefile – Crypto‑Exporter
# ─────────────────────────────────────────────────────────────────────────────
#  Variáveis de configuração (pode ser sobrescritas na linha de comando)
# ─────────────────────────────────────────────────────────────────────────────
APP_NAME      := crypto-exporter
PKG_ROOT      := ./cmd/$(APP_NAME)
BUILD_DIR     := ./bin
VERSION       ?= $(shell git describe --tags --always --dirty)
COMMIT        ?= $(shell git rev-parse --short HEAD)
BUILD_TIME    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS       := -X main.version=$(VERSION) \
                 -X main.commit=$(COMMIT) \
                 -X main.buildTime=$(BUILD_TIME)

#  Ferramentas externas (instaláveis via go install)
GOIMPORTS     := $(shell go env GOPATH)/bin/goimports
GOLINT        := $(shell go env GOPATH)/bin/golangci-lint
MOCKGEN       := $(shell go env GOPATH)/bin/mockgen
GOVULNCHECK   := $(shell go env GOPATH)/bin/govulncheck

#  Targets padrão
.PHONY: all
all: build ## (default) Compila o binário e executa os testes

# ─────────────────────────────────────────────────────────────────────────────
#  Compilação
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: build
build: $(BUILD_DIR)/$(APP_NAME) ## Compila o binário na pasta ./bin

$(BUILD_DIR)/$(APP_NAME): $(shell find . -name '*.go' -not -path "./vendor/*")
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $@ $(PKG_ROOT)

# ─────────────────────────────────────────────────────────────────────────────
#  Testes
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: test
test: ## Executa todos os testes unitários (race + cobertura)
	go test ./... -race -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -n1

.PHONY: test-ci
test-ci: ## Versão “CI” – falha se cobertura < 80%
	go test ./... -race -coverprofile=coverage.out
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' | \
	    awk -F. '{if ($$1 < 80) {exit 1}}'

# ─────────────────────────────────────────────────────────────────────────────
#  Lint / Formatação
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: fmt
fmt: ## Formata o código (go fmt + goimports)
	go fmt ./...
	$(GOIMPORTS) -w .

.PHONY: lint
lint: ## Executa o golangci‑lint (deve estar instalado)
	$(GOLINT) run ./...

# ─────────────────────────────────────────────────────────────────────────────
#  Segurança
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: vuln
vuln: ## Roda govulncheck (detecta vulnerabilidades conhecidas)
	$(GOVULNCHECK) ./...

# ─────────────────────────────────────────────────────────────────────────────
#  Mocks (para testes)
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: mocks
mocks: ## (Re)gera os mocks a partir das interfaces
	$(MOCKGEN) -source=internal/client/coingecko.go -destination=internal/client/mock_coingecko.go -package=client

# ─────────────────────────────────────────────────────────────────────────────
#  Execução
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: run
run: build ## Executa o binário recém‑compilado
	$(BUILD_DIR)/$(APP_NAME)

.PHONY: dev
dev: ## Executa com hot‑reload (air, reflex ou entr – escolha o que prefere)
	# Exemplo usando "air" (go get -u github.com/cosmtrek/air)
	air

# ─────────────────────────────────────────────────────────────────────────────
#  Docker
# ─────────────────────────────────────────────────────────────────────────────
DOCKER_REGISTRY ?= ghcr.io/your-org
DOCKER_IMAGE    := $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)

.PHONY: docker-build
docker-build: ## Constrói a imagem Docker (multi‑stage)
	@docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-push
docker-push: docker-build ## Faz push da imagem para o registry configurado
	@docker push $(DOCKER_IMAGE)

.PHONY: docker-run
docker-run: ## Roda a imagem localmente (exemplo: mapeia a porta 8081)
	@docker run --rm -p 8081:8081 $(DOCKER_IMAGE)

# ─────────────────────────────────────────────────────────────────────────────
#  Limpeza
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: clean
clean: ## Remove binários, arquivos de cobertura e caches
	@rm -rf $(BUILD_DIR) coverage.out
	@go clean -testcache -modcache

# ─────────────────────────────────────────────────────────────────────────────
#  Ajuda
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: help
help:
	@printf "\n\033[1;34mMakefile targets for $(APP_NAME)\033[0m\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | \
	    sort | awk 'BEGIN {FS = ":.*?## "}; \
	    {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ─────────────────────────────────────────────────────────────────────────────
#  Definições implícitas
# ─────────────────────────────────────────────────────────────────────────────
# Se alguma das ferramentas externas não existir, tenta instalá‑las.
$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@latest

$(GOLINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(MOCKGEN):
	go install github.com/golang/mock/mockgen@latest

$(GOVULNCHECK):
	go install golang.org/x/vuln/cmd/govulncheck@latest

