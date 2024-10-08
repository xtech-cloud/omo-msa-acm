APP_NAME := xm-msa-acm
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
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.IsPermission '{"user":"Tom4", "path":"omo/msa/acm/menu/remove", "action":"get"}'

.PHONY: tester
tester:
	go build -o ./bin/ ./tester

.PHONY: dist
dist:
	mkdir -p dist
	rm -f dist/${APP_NAME}-${BUILD_VERSION}.tar.gz
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build . -t omo.msa.acm:latest

.PHONY: updev
updev:
	scp -P 9700 dist/${APP_NAME}-${BUILD_VERSION}.tar.gz root@192.168.1.10:/root/

.PHONY: upload2
upload2:
	scp -P 9099 dist/${APP_NAME}-${BUILD_VERSION}.tar.gz root@47.93.209.105:/root/

.PHONY: upload
upload:
	scp -P 9099 dist/${APP_NAME}-${BUILD_VERSION}.tar.gz root@101.200.166.80:/root/