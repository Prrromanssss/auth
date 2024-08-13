LOCAL_BIN:=$(CURDIR)/app/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

generate:
	make generate-user-api

generate-user-api:
	mkdir -p app/pkg/user_v1
	protoc --proto_path app/api/user_v1 \
	--proto_path app/vendor.protogen  \
	--go_out=app/pkg/user_v1 \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=app/bin/protoc-gen-go \
	--go-grpc_out=app/pkg/user_v1 \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=app/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=app/pkg/user_v1 \
	--grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=app/bin/protoc-gen-grpc-gateway \
	--validate_out lang=go:app/pkg/user_v1 \
	--validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=app/bin/protoc-gen-validate \
	app/api/user_v1/user.proto

local-migration-status:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DIR} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DIR} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DIR} down -v

test-coverage:
	@cd app && \
	go clean -testcache && \
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Prrromanssss/auth/internal/service/...,github.com/Prrromanssss/auth/internal/api/... -count 5  && \
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out  && \
	rm coverage.tmp.out
	cd app && go tool cover -html=coverage.out;
	cd app &&go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
	@if [ ! -d app/vendor.protogen/google ]; then \
	git clone https://github.com/googleapis/googleapis app/vendor.protogen/googleapis &&\
	mkdir -p  app/vendor.protogen/google/ &&\
	mv app/vendor.protogen/googleapis/google/api app/vendor.protogen/google &&\
	rm -rf app/vendor.protogen/googleapis ;\
	fi
	@if [ ! -d app/vendor.protogen/validate ]; then \
		mkdir -p app/vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate app/vendor.protogen/protoc-gen-validate &&\
		mv app/vendor.protogen/protoc-gen-validate/validate/*.proto app/vendor.protogen/validate &&\
		rm -rf app/vendor.protogen/protoc-gen-validate ;\
	fi

run-local:
	docker-compose -f .docker-compose.yaml up --build -d

down-local:
	docker-compose -f .docker-compose.yaml down