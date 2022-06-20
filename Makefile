BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY} ./cmd/server/*.go

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t order-service .

run:
	docker-compose up -d

stop:
	docker-compose down

lint:
	golangci-lint run 

sqlc:
	sqlc generate

.PHONY: test engine unittest clean docker run stop lint sqlc