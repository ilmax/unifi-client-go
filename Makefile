.PHONY: generate-openapi

OPENAPI_BASE_URL ?= https://developer.ui.com/network
API_VERSION ?=

# Generate OpenAPI spec, patched codegen input, and client code.
generate-openapi:
	@if [ -z "$(API_VERSION)" ]; then \
		echo "Error: API_VERSION is required (e.g., 10.4.57)"; \
		exit 1; \
	fi
	@normalized="$${API_VERSION#v}"; \
	spec_dir="./openapi/unifi-network"; \
	spec_path="$$spec_dir/$$normalized.json"; \
	mkdir -p "$$spec_dir"; \
		curl -fsSL "$(OPENAPI_BASE_URL)/v$$normalized/openapi.json" -o "$$spec_path" || { \
			echo "Error: Unable to download OpenAPI spec for version $$normalized"; \
			exit 1; \
		}; \
	go run ./cmd/openapi-preprocess \
		-in "$$spec_path" \
		-out "$$spec_dir/$$normalized.codegen.json"; \
	find ./pkg/network -maxdepth 1 -type f -name '*.go' ! -name 'network.go' -delete; \
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0 \
		-generate types,client \
		-package network \
		-o ./pkg/network/openapi.gen.go \
		"$$spec_dir/$$normalized.codegen.json"
