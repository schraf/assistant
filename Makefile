BUILD_DIR=bin
REGION ?= us-central1

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/tokens ./cmd/tokens
	go build -o $(BUILD_DIR)/services ./cmd/service
	go build -o $(BUILD_DIR)/job ./cmd/job

generate-token:
	@echo "Generating telegraph token..."
	go run ./cmd/tokens

vet:
	@echo "Vetting code..."
	go vet ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

test:
	@echo "Running tests..."
	go test ./internal/test/...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

terraform-init:
	@echo "Initializing Terraform..."
	cd terraform && terraform init

terraform-refresh:
	@echo "Refreshing Terraform state..."
	cd terraform && terraform refresh

terraform-plan:
	@echo "Planning Terraform changes..."
	cd terraform && terraform plan

terraform-apply:
	@echo "Applying Terraform changes..."
	cd terraform && terraform apply

terraform-bootstrap:
	@echo "Applying Terraform changes with initial placeholder image..."
	@echo "This is for the first apply only, to set up Cloud Build infrastructure."
	@echo "After this, use 'make terraform-apply' for subsequent applies."
	cd terraform && terraform apply -var="initial_image=gcr.io/cloudrun/hello"

terraform-destroy:
	@echo "Destroying Terraform resources..."
	cd terraform && terraform destroy

terraform-validate:
	@echo "Validating Terraform configuration..."
	cd terraform && terraform validate

terraform-fmt:
	@echo "Formatting Terraform files..."
	cd terraform && terraform fmt -recursive

deploy:
	@echo "Building and deploying with Cloud Build..."
	@if [ -z "$(PROJECT_ID)" ]; then \
		echo "Error: PROJECT_ID environment variable is required"; \
		echo "Usage: PROJECT_ID=your-project-id make deploy"; \
		exit 1; \
	fi
	@TAG=$$(git rev-parse --short HEAD 2>/dev/null); \
	if [ -z "$$TAG" ]; then \
		TAG=$$(date +%s); \
	fi; \
	echo "Using tag: $$TAG"; \
	gcloud builds submit --config=cloudbuild.yaml \
		--substitutions=_REGION=$(REGION),_REPO_NAME=assistant,_TAG=$$TAG \
		--project=$(PROJECT_ID)

help:
	@echo "Available targets:"
	@echo "  all                 - Run vet and build"
	@echo "  build               - Build the application"
	@echo "  run                 - Run the application"
	@echo "  generate-token      - Generate a Telegraph API token"
	@echo "  vet                 - Vet the code"
	@echo "  fmt                 - Format the code"
	@echo "  test                - Run tests"
	@echo "  clean               - Clean build artifacts"
	@echo "  deps                - Install dependencies"
	@echo ""
	@echo "Terraform targets:"
	@echo "  terraform-init      - Initialize Terraform"
	@echo "  terraform-refresh   - Refreshes the Terraform state"
	@echo "  terraform-plan      - Plan Terraform changes"
	@echo "  terraform-apply     - Apply Terraform changes (use after initial setup)"
	@echo "  terraform-bootstrap - Apply Terraform with placeholder image (first apply only)"
	@echo "  terraform-destroy   - Destroy Terraform resources"
	@echo "  terraform-validate  - Validate Terraform configuration"
	@echo "  terraform-fmt       - Format Terraform files"
	@echo ""
	@echo "Deployment targets:"
	@echo "  deploy              - Build and deploy using Cloud Build (requires PROJECT_ID)"
	@echo ""
	@echo "  help                - Show this help message"
