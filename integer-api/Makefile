dep:
	go get -t ./...

build:
	go build -o integer-api *.go

build-static:
	CGO_ENABLED=0 go build -v -a -installsuffix cgo -o integer-api *.go

docker-build:
	docker build -f Dockerfile.build -t loopstock.io/integer-api:build .
	docker create --name integer-api-builded loopstock.io/integer-api:build
	docker cp integer-api-builded:/go/src/github.com/guilherme-santos/loopstock.io/integer-api/integer-api .
	docker rm -f integer-api-builded
	docker image rm loopstock.io/integer-api:build

docker-image: docker-build
	docker build -t loopstock.io/integer-api .

run:
	docker-compose up -d

logs:
	docker-compose logs -f ${SERVICE}

stop:
	docker-compose stop

cleanup:
	docker-compose down -v

integer-api-port:
	@docker-compose port integer-api 8080
