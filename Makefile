SHELL=/bin/bash
ROOT_DIR := $(shell pwd)
REGISTRY=register.docker.com
IMAGE_TAG := $(shell git rev-parse HEAD)
IMAGE_NAME := auth2-jwt

.PHONY: ci
ci: lint svc test

.PHONY: mod
mod:
	go mod download
	go mod vendor

.PHONY: grpcgen
grpcgen:
	go build -o protoc-gen-go ./vendor/github.com/golang/protobuf/protoc-gen-go
	protoc -I ./vendor -I api api/service.proto --plugin=./protoc-gen-go --go_out=plugins=grpc:api
	go build ./api

.PHONY: svc
svc:
	go build -mod=vendor -o artifacts/svc ./cmd/svc

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -v `go list ./...`

.PHONY: mockgen
mock_service_storage:
	mockgen -source=srv/http/deps.go -destination=service/mock/deps.go

.PHONY: dockerise
dockerise:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -f Dockerfile .
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: push_image_to_registry
push_image_to_registry:
	`AWS_SHARED_CREDENTIALS_FILE=~/.aws/credentials AWS_PROFILE=prof_name aws ecr get-login --region us-west-2 --no-include-email`
	docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
	#docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}/${IMAGE_NAME}:latest
	#docker push ${REGISTRY}/${IMAGE_NAME}:latest
