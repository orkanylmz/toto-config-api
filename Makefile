include .env
export

.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
		@./scripts/openapi-http.sh skuconfig internal/skuconfig/ports ports

.PHONY: lint
lint:
		@go-cleanarch
		@./scripts/lint.sh common
		@./scripts/lint.sh skuconfig

.PHONY: fmt
fmt:
		goimports -l -w internal/

test:
		@./scripts/test.sh common .e2e.env
		@./scripts/test.sh skuconfig .test.env