BIN_PATH = build/bin
BUILD_CMD = GOOS=linux GOARCH=amd64 go build -v -ldflags '-s -w' -a -tags netgo -installsuffix netgo -o ${BIN_PATH}/app
.PHONY: build
build-local:
	${BUILD_CMD} ./cmd/payment/http/main.go

build-lambda:
	${BUILD_CMD} ./cmd/payment/aws/main.go