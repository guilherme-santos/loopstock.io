dep:
	go get -t ./...

build:
	go build -o integeraverage-cal main.go

build-static:
	CGO_ENABLED=0 go build -v -a -installsuffix cgo -o integeraverage-cal main.go

docker-build:
	docker build -f Dockerfile.build -t loopstock.io/integeraverage-cal:build .
	docker create --name integeraverage-cal-builded loopstock.io/integeraverage-cal:build
	docker cp integeraverage-cal-builded:/go/src/github.com/guilherme-santos/loopstock.io/integeraverage-cal/integeraverage-cal .
	docker rm -f integeraverage-cal-builded
	docker image rm loopstock.io/integeraverage-cal:build

docker-image: docker-build
	docker build -t loopstock.io/integeraverage-cal .

run:
	docker-compose up -d

logs:
	docker-compose logs -f ${SERVICE}

stop:
	docker-compose stop

cleanup:
	docker-compose down -v
