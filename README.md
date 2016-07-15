# Docker Dispatch
A service that dispatches Docker containers from RabbitMQ messages.


## Architecture

Docker dispatch listens for properly formatted messages from RabbitMQ and creates/runs containers from pre-built images in Docker. It talks to the Docker daemon via its RESTful API and delivers responses in stdout.
```
rabbitmq ----(message)----> docker-dispatch <----(REST Api)----> dockerd
```

## Installation

With go: `go get github.com/mutaphore/docker-dispatch`


## Running

To run the dispatcher service: 
```
$ ./docker-dispatch [options] dockerHostAddr amqpAddr
```

`dockerHostAddr` is the docker daemon host address you would like the dispatcher to connect to. This can be expressed as a TCP url or the socket path. A valid TCP url would have **ip:port** format such as `172.17.0.1:2375`. A socket path is the unix file path to the socket which is usually like `/var/run/docker.sock`. Information on how to configure and run the Docker daemon can be found [here](https://docs.docker.com/engine/admin/configuring/)

`amqpAddr` is the full url to RabbitMQ with format **amqp://username:password@host:port**. Here's the full [spec](https://www.rabbitmq.com/uri-spec.html). For example, a valid url would be something like `amqp://guest:guest@localhost:5672/`

Example:
```
./dockerdispatch -q myqueue 172.17.0.1:2375 amqp://guest:guest@localhost:5672
```


## Options
```
Options:
-q, queue name
-v, verbose
```
