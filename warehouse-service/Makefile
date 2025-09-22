export LINTER_VERSION ?= 1.55.2

GO_PACKAGES ?= $(shell go list ./... | grep -v 'examples\|qtest\|mock\|config\|cmd\|dev\|testutil\|/server')
MODULES      = $(shell cd module && ls -d */)
TMP_DIR     := $(shell mktemp -d)
UNAME       := $(shell uname)

ODIR 		:= build/_output
DOCKERFILES  = $(shell cd build && find */ -name 'Dockerfile' -print)

# Default mysql migration settings
export MYSQL_USER ?= root
export MYSQL_PASS ?= root
export MYSQL_HOST ?= localhost
export MYSQL_PORT ?= 3306
export MYSQL_DB   ?= warehouse_service

export VAR_SERVICES ?= $(foreach path, $(DOCKERFILES), $(path:%/Dockerfile=%))

dep:
	@go mod tidy

bin:
	@mkdir -p bin



tool-migrate: bin
ifneq (,$(wildcard bin/migrate))
    # do not download again
else ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo "Your OS is not supported."
endif

migrate: tool-migrate
	@$(foreach module, $(MODULES), cp module/$(module)/db/migrate/*.sql $(TMP_DIR) 2>/dev/null;)
	@bin/migrate -source file://$(TMP_DIR) -database "mysql://$(MYSQL_USER):$(MYSQL_PASS)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DB)" $(MIGRATE_ARGS) $(N)

generate-db-migration: tool-migrate
	@bin/migrate create -ext sql -dir module/$(MIGRATE_MODULE)/db/migrate $(MIGRATE_NAME)



tool-mgen:
	@go install github.com/golang/mock/mockgen@v1.6.0

mgen: tool-mgen
	@go generate ./...


test:
	@go test -race -v ${GO_PACKAGES}

coverage:
	@go run gotest.tools/gotestsum@latest --format pkgname --jsonfile UT-warehouse-service.json -- -coverprofile=coverage.out ${GO_PACKAGES}
	@go tool cover -func=coverage.out


$(ODIR):
	@mkdir -p $(ODIR)

compile: $(ODIR)
	@$(foreach svc, $(VAR_SERVICES), \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(ODIR)/$(svc) cmd/$(svc)/main.go;)

build: $(ODIR)
	@$(foreach svc, $(VAR_SERVICES), \
		docker build -t commerce-exercise-warehouse-service/$(svc) -f ./build/$(svc)/Dockerfile .;)

