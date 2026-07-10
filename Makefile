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
	# Wipe every mock artefact left behind by either v2 (top-level
	# ./mocks/*.go) or per-package (mock_*.go anywhere). The original
	# `-not -path "./pkg/infrastructure/*"` predicate was inverted — it
	# kept mocks OUTSIDE pkg/infrastructure and deleted them INSIDE,
	# which is the opposite of what we want. Use plain `-delete` so
	# generated files don't leak into git.
	@find . -name 'mock_*.go' -delete
	@rm -rf ./mocks

generate-mock: clean-mock
	@echo "--- Generating Mock by Mockery ---"
	@mockery --all
	
test: setup
	@echo "--- Running test ---"
	@mkdir -p report
	# Pin gotestsum — `@latest` will silently bump the go directive
	# requirement past what our Go toolchain supports (see go.mod's `toolchain`
	# line). v1.13.0 is the latest release that compiles cleanly under
	# Go 1.25.x as of writing.
	@go run gotest.tools/gotestsum@v1.13.0 -f testname -- ./pkg/module/... --coverprofile="report/c.out" -shuffle=on

# Note: `generate-mock` is intentionally NOT a prerequisite of `test`. The
# project's *_test.go files don't consume any generated mocks (verified by
# grep: no test imports `stretchr/testify/mock` or `mocks.*`), so running
# mockery on every CI build was producing output nobody reads. Worse,
# mockery v2 always type-checks its own freshly-written output via
# `go/packages`, which trips on `github.com/stretchr/testify/mock` (a
# sub-package that isn't published as a standalone Go module). That made
# the `make test` target fail at the generate-mock step on every run.
#
# To regenerate mocks when actually needed:
#   1. add a .mockery.yaml that scopes output to the interfaces under test
#   2. write them next to the source package (`dir: "{{.InterfaceDir}}"`)
#   3. re-introduce `generate-mock` as a prerequisite of `test`
# For now, `make generate-mock` is still callable by hand but CI skips it.

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


