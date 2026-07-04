.PHONY: generate-mock docs

setup:
	@echo "--- Setup and generated config yaml files ---"
	@mkdir -p config/resources
	@cp -r config/example/*.yml config/resources/

api:
	@echo "--- Running api server in dev mode ---"
	@go run main.go api

docs:
	@echo "--- Generating swagger docs ---"
	@swag fmt
	@swag init --parseInternal --parseVendor --parseDependencyLevel 3

clean-mock:
	@echo "--- Cleaning Mock ---"
	@find . -name 'mock_*.go' -not -path "./pkg/infrastructure/*" -delete

generate-mock: clean-mock
	@echo "--- Generating Mock by Mockery ---"
	@mockery --all
	
test: generate-mock
	@echo "--- Running test ---"
	@mkdir -p report
	@go run gotest.tools/gotestsum@latest -f testname -- ./pkg/module/... ./pkg/clients/... --coverprofile="report/c.out" -shuffle=on

test-coverage: test
	@echo "--- Running test coverage ---"
	@grep -v -e "mock_" -e "test_spec" -e "/pkg/clients/.*/entity/.*" report/c.out > report/coverage.out
	@go tool cover -html="report/coverage.out" -o "report/coverage.html"

generate-password:
	@echo "--- Generating password ---"
	@go run main.go password generate $(password)

migrate-up:
	@echo "--- Running db migration up ---"
	@go run main.go db migrate up $(step) --schema=$(schema)

migrate-down:
	@echo "--- Running db migration down ---"
	@go run main.go db migrate down $(step) --schema=$(schema)

migrate-create:
	@echo "--- Creating db migration ---"
	@migrate create -ext sql -dir db/migrations/postgres -seq $(name)

build: setup
	@echo "--- Building binary file ---"
	@go build -o ./main main.go


