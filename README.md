# Docker Dispatch
A service that dispatches Docker containers from RabbitMQ messages.


## Installation

With go: `go get -u github.com/mutaphore/docker-dispatch`


## Running

To run the dispatcher service: `$ ./docker-dispatch [options] dockerIpAndPort amqpUrl`

`dockerIpAndPort` is the docker daemon host address you would like the dispatcher to connect to. This can be expressed as a TCP url or the socket path. A valid TCP url would have **ip:port** format such as `172.17.0.1:2375`. A socket path is something like `/var/run/docker.sock`