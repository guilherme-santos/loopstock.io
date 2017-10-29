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
curl localhost:<api-port>/v1/integers/last
```

## Question

### Describe how you would provision a set of linux machines and how'd you orchestrate, monitor and troubleshoot these services?

This project generate 3 docker images `loopstock.io/integer-gen`, `loopstock.io/integeraverage-cal` and `loopstock.io/integer-api` and they're configured using environment variables, it means their deploy is really easy using AWS ECS ([deploy-ecs](https://github.com/guilherme-santos/deploy-ecs) is a tool that I created to make easier you build your image, send it to ECR, update your task-definition and do the deploy itself). Another option is use GCP - Google Cloud Platform with Kubernetes, the third option I'd say to install and configure Kubernetes by your own.

So, in general I think you can work with these 3 possibilities: AWS, GCP or building everything by your own. To deploy and orchestrate AWS ECS, GPC Kubernetes or just Kubernetes. To monitor, AWS and GPC will give to you really nice tools with graphs etc. but if you decide do it everything you can use tools like InfluxDB, Statsd, and Grafana or even some others payed tools. For troubleshoot both platform have good tools (like AWS CloudWatch), but you can also use ELK Stack.

To database and Broker instances, you can use AWS RDS and AWS SQS, or use **Container Linux Configuration** from CoreOS to install and configure your servers.

## Bonus point

Related to the bonus point I thought a little bit about the problem, but I didn't really implement anything because I think the problem is not so clear to me, I mean I cannot see it as a real problem. But some options...

- I could have a main generator that is shared between all **integer-gen** instances, the problem this instance cannot be scaled.
- I could make the random seed configurable, it means you can set the seed to be used, the problem is each instance need to have a different seed what I believe is difficult to manage.
- I could have a new service in front of **integer-gen** that discard in case both services generate the same number, it looks for me a big overhead .
- I could also check on **integeraverage-cal** if there is a duplicated number, maybe I'll need to identify which instance generate the number, but I don't think it's a big problem.

That was some ways I thought to solve, I didn't really like so much all these solution, that's why I didn't implement and I'd prefer understand a little bit more what's the impact of generate the same number at the same time.
