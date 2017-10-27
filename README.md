# Loopstock.io Challenge

## Before start

First of all you need build all images, to do that, simply type:
```
make docker-image
```

## Running

To run all services just type:
```
make run
```

To stop, you can use:
```
make stop
```

To cleanup all docker containers (including networks and volumes). P.S.: *docker images will not be deleted*
```
make cleanup
```

#### Ports

Some services are exposing ports, to check these ports type:
```
make ports
```

#### NSQAdmin - Message Broker

To monitor NSQ access: `http://localhost:<nsqadmin-port>`

#### API

To get last 10 averages, type:
```
curl localhost:<api-port>/v1/integers
```

To get just the last message, type:
```
curl localhost:8080/v1/integers/last
```

## Question

### Describe how you would provision a set of linux machines and how'd you orchestrate, monitor and troubleshoot these services?

##### Provision and Orchestration

This project generate 3 docker images `loopstock.io/integer-gen`, `loopstock.io/integeraverage-cal` and `loopstock.io/integer-api` and they're configured using envvars it means that their deploy is really easy using AWS ECS ((deploy-ecs)<https://github.com/guilherme-santos/deploy-ecs> is a tool that I created to make easier you build your image, send it to ECR, update your task-definition and do the deploy itself), if it's not an option use AWS ECS I'd use Kubernets and believe I'll have the same result.

For the other servers MySQL, NSQd, I'd use **Container Linux Configuration** from CoreOS.
