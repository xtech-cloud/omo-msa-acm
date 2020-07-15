APP_NAME := omo.msa.acm
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=. grpc/proto/common.proto
	protoc --proto_path=. --micro_out=. --go_out=. grpc/proto/user.proto
	protoc --proto_path=. --micro_out=. --go_out=. grpc/proto/role.proto
	protoc --proto_path=. --micro_out=. --go_out=. grpc/proto/menu.proto

.PHONY: build
build:
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: call
call:
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.RemoveOne '{"uid":"John"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.GetList '{"page":3, "number":5}'

.PHONY: tester
tester:
	go build -o ./bin/ ./tester

.PHONY: dist
dist:
	mkdir dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo.msa.acm:latest
