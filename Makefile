MAIN_PATH=cmd/blcheck/main.go
BIN_PATH=bin/blcheck

build:
	go build -o ${BIN_PATH} ${MAIN_PATH}

build-all:
	GOARCH=amd64 GOOS=darwin go build -o ${BIN_PATH}-darwin ${MAIN_PATH}
	GOARCH=amd64 GOOS=linux go build -o ${BIN_PATH}-linux ${MAIN_PATH}
	GOARCH=amd64 GOOS=windows go build -o ${BIN_PATH}-windows ${MAIN_PATH}

clean:
	go clean
	rm -rf ${BIN_PATH}*
	

run-build:
	./${BIN_PATH} 

run:
	go run ${MAIN_PATH}

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet ./...

# TODO: Add linting and checkout tool
# lint:
# 	golangci-lint run --enable-all