INSTANCES ?= 2

dep:
	go get -t ./...

build:
	go build -o integer-gen main.go

build-static:
	CGO_ENABLED=0 go build -v -a -installsuffix cgo -o integer-gen main.go

docker-build:
	docker build -f Dockerfile.build -t loopstock.io/integer-gen:build .
	docker create --name integer-gen-builded loopstock.io/integer-gen:build
	docker cp integer-gen-builded:/go/src/github.com/guilherme-santos/loopstock.io/integer-gen/integer-gen .
	docker rm -f integer-gen-builded
	docker image rm loopstock.io/integer-gen:build

docker-image: docker-build
	docker build -t loopstock.io/integer-gen .

run:
	docker-compose up --scale integer-gen=${INSTANCES} -d

logs:
	docker-compose logs -f

stop:
	docker-compose stop

cleanup:
	docker-compose down -v
