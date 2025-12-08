.PHONY: help test test-unit test-integration test-e2e test-all build clean deploy-aws deploy-azure terraform-aws terraform-azure golangci-lint golangci-lint-fix install-golangci-lint

# Variables
GO := go
DOCKER_COMPOSE := docker compose
AWS_TERRAFORM_DIR := src/infrastructure/clouds/aws/terraform
AZURE_TERRAFORM_DIR := src/infrastructure/clouds/azure/terraform
AWS_FUNCTIONS_DIR := src/infrastructure/clouds/aws/functions
AZURE_FUNCTIONS_DIR := src/infrastructure/clouds/azure/functions

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

##@ General

help: ## Show this help message
	@echo "$(CYAN)Available commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

##@ Testing

test: test-all ## Run all tests (unit + integration + e2e)

test-unit: ## Run unit tests
	@echo "$(CYAN)🧪 Running unit tests...$(NC)"
	$(GO) test ./src/... -v -cover

test-integration: ## Run integration tests with Docker
	@echo "$(CYAN)🧪 Running integration tests...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.test.yml up --abort-on-container-exit --exit-code-from goprojectskeleton-integration-tests --build

test-integration-local: ## Run integration tests locally (without Docker)
	@echo "$(CYAN)🧪 Running integration tests locally...$(NC)"
	$(GO) test ./tests/integration/... -v -count=1 -cover

test-e2e: ## Run E2E tests with Docker (Bruno)
	@echo "$(CYAN)🧪 Running E2E tests...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.e2e.yml up --abort-on-container-exit --build

test-all: test-unit test-integration test-e2e ## Run all test types

test-clean: ## Clean test containers and volumes
	@echo "$(CYAN)🧹 Cleaning test containers...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.test.yml down -v
	$(DOCKER_COMPOSE) -f docker/docker-compose.e2e.yml down -v

##@ Build

build: ## Build the project
	@echo "$(CYAN)🔨 Building project...$(NC)"
	$(GO) build -o bin/goprojectskeleton ./src/infrastructure/server/cmd

build-aws-functions: ## Generate AWS Lambda functions
	@echo "$(CYAN)🔨 Generating AWS Lambda functions...$(NC)"
	cd $(AWS_FUNCTIONS_DIR) && $(GO) run main.go generate

build-azure-functions: ## Generate Azure Functions
	@echo "$(CYAN)🔨 Generating Azure Functions...$(NC)"
	cd $(AZURE_FUNCTIONS_DIR) && $(GO) run generate.go functions.go

##@ AWS Deployment

deploy-aws: build-aws-functions ## Deploy all AWS Lambda functions
	@echo "$(CYAN)☁️  Deploying AWS Lambda functions...$(NC)"
	cd $(AWS_FUNCTIONS_DIR) && $(GO) run main.go deploy

deploy-aws-function: ## Deploy a specific AWS Lambda function (usage: make deploy-aws-function FUNCTION=health-check)
	@if [ -z "$(FUNCTION)" ]; then \
		echo "$(RED)❌ Error: You must specify FUNCTION=function-name$(NC)"; \
		echo "Example: make deploy-aws-function FUNCTION=health-check"; \
		exit 1; \
	fi
	@echo "$(CYAN)☁️  Deploying AWS Lambda function: $(FUNCTION)...$(NC)"
	cd $(AWS_FUNCTIONS_DIR) && $(GO) run main.go deploy $(FUNCTION)

terraform-aws-init: ## Initialize Terraform for AWS
	@echo "$(CYAN)🔧 Initializing Terraform for AWS...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform init

terraform-aws-plan: terraform-aws-init ## Show Terraform plan for AWS
	@echo "$(CYAN)📋 Generating Terraform plan for AWS...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform plan

terraform-aws-apply: terraform-aws-init ## Apply Terraform changes for AWS
	@echo "$(CYAN)🚀 Applying Terraform changes for AWS...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform apply

terraform-aws-destroy: ## Destroy AWS infrastructure
	@echo "$(YELLOW)⚠️  Destroying AWS infrastructure...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform destroy

terraform-aws-output: ## Show Terraform outputs for AWS
	@echo "$(CYAN)📤 Showing Terraform outputs for AWS...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform output

terraform-aws-validate: ## Validate Terraform configuration for AWS
	@echo "$(CYAN)✅ Validating Terraform configuration for AWS...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform validate

terraform-aws-fmt: ## Format Terraform files for AWS
	@echo "$(CYAN)📝 Formatting Terraform files for AWS...$(NC)"
	cd $(AWS_TERRAFORM_DIR) && terraform fmt -recursive

##@ Azure Deployment

deploy-azure: build-azure-functions ## Deploy Azure Functions
	@echo "$(CYAN)☁️  Deploying Azure Functions...$(NC)"
	@echo "$(YELLOW)⚠️  Note: Azure Functions deployment is done via Terraform or Azure CLI$(NC)"
	@echo "$(YELLOW)    Use 'make terraform-azure-apply' to deploy the infrastructure$(NC)"

terraform-azure-init: ## Initialize Terraform for Azure
	@echo "$(CYAN)🔧 Initializing Terraform for Azure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform init

terraform-azure-plan: terraform-azure-init ## Show Terraform plan for Azure
	@echo "$(CYAN)📋 Generating Terraform plan for Azure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform plan

terraform-azure-apply: terraform-azure-init ## Apply Terraform changes for Azure
	@echo "$(CYAN)🚀 Applying Terraform changes for Azure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform apply

terraform-azure-destroy: ## Destroy Azure infrastructure
	@echo "$(YELLOW)⚠️  Destroying Azure infrastructure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform destroy

terraform-azure-output: ## Show Terraform outputs for Azure
	@echo "$(CYAN)📤 Showing Terraform outputs for Azure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform output

terraform-azure-validate: ## Validate Terraform configuration for Azure
	@echo "$(CYAN)✅ Validating Terraform configuration for Azure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform validate

terraform-azure-fmt: ## Format Terraform files for Azure
	@echo "$(CYAN)📝 Formatting Terraform files for Azure...$(NC)"
	cd $(AZURE_TERRAFORM_DIR) && terraform fmt -recursive

##@ Utilities

clean: ## Clean generated files and binaries
	@echo "$(CYAN)🧹 Cleaning generated files...$(NC)"
	rm -rf bin/
	rm -rf tmp/
	$(GO) clean -cache
	$(GO) clean -modcache

clean-docker: ## Clean Docker containers, images and volumes
	@echo "$(CYAN)🧹 Cleaning Docker...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.dev.yml down -v
	$(DOCKER_COMPOSE) -f docker/docker-compose.test.yml down -v
	$(DOCKER_COMPOSE) -f docker/docker-compose.e2e.yml down -v
	docker system prune -f

deps: ## Download and install dependencies
	@echo "$(CYAN)📦 Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

deps-update: ## Update dependencies
	@echo "$(CYAN)📦 Updating dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy

fmt: ## Format Go code
	@echo "$(CYAN)📝 Formatting Go code...$(NC)"
	$(GO) fmt ./...

vet: ## Run go vet on the code
	@echo "$(CYAN)🔍 Running go vet...$(NC)"
	$(GO) vet ./...

golangci-lint: ## Run golangci-lint
	@echo "$(CYAN)🔍 Running golangci-lint...$(NC)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(RED)❌ Error: golangci-lint is not installed$(NC)"; \
		echo "Install it with: make install-golangci-lint"; \
		exit 1; \
	fi
	golangci-lint run ./...

golangci-lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(CYAN)🔍 Running golangci-lint with auto-fix...$(NC)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(RED)❌ Error: golangci-lint is not installed$(NC)"; \
		echo "Install it with: make install-golangci-lint"; \
		exit 1; \
	fi
	golangci-lint run --fix ./...

install-golangci-lint: ## Install golangci-lint
	@echo "$(CYAN)📦 Installing golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠️  golangci-lint is already installed$(NC)"; \
	else \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
		echo "$(GREEN)✅ golangci-lint installed successfully$(NC)"; \
	fi

lint: fmt vet golangci-lint ## Run code formatting and validation

dev: ## Start development environment with Docker
	@echo "$(CYAN)🚀 Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.dev.yml up -d

dev-down: ## Stop development environment
	@echo "$(CYAN)🛑 Stopping development environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.dev.yml down

dev-logs: ## Show development environment logs
	@echo "$(CYAN)📋 Showing development environment logs...$(NC)"
	$(DOCKER_COMPOSE) -f docker/docker-compose.dev.yml logs -f
