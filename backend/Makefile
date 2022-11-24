BINARY_NAME=route_planner

build:
	go build -o ${BINARY_NAME} cmd/main.go

run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}