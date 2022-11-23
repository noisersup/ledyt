BIN_NAME=bin/ledyt
CMD_LOCATION=cmd/ledyt

all: test build 

run:
	go run ${CMD_LOCATION}/main.go

build:
	go build -o ${BIN_NAME} ${CMD_LOCATION}/main.go

test:
	go test -v -race ./...

clean:
	go clean
