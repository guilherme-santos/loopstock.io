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

Also to cleanup all docker containers (including networks and volumes). P.S.: *Images will not be deleted*
```
make cleanup
```

## Question

### Describe how you would provision a set of linux machines and how'd you orchestrate, monitor and troubleshoot these services?