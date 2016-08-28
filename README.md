# Docker Dispatch

## Architecture

Docker dispatch listens for properly formatted messages from RabbitMQ and creates/runs containers from pre-built images in Docker. It talks to the Docker daemon via its RESTful API and delivers responses in stdout.

![architecture](res/docker-dispatch-architecture.png)

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

For more information on how to setup Docker to bind to different addresses see [here](https://docs.docker.com/engine/reference/commandline/dockerd/#bind-docker-to-another-host-port-or-a-unix-socket).


## Messages formats

Docker dispatcher accepts JSON messages similar to the command that you would give directly to the docker daemon. Note that all key values are capitalized.

### Running a container from an image (docker run)

* Dockercmd - docker command to execute. Currently supports run, stop and remove.
* Options - options to pass into the docker command.
** Name - name of container
** Entrypoint - command to execute in container
** Attach - array of strings selected from "STDIN", "STDOUT", "STDERR"
** Remove - boolean, whether or not to remove container after exit
* Image - image name to create container from
* Cmd - exec form of command to run in container having format ["executable","param1","param2"]

#### Example
```javascript
{
  Dockercmd: "run",
  Options: {
    Name: "hello_world",
    Attach: ["STDOUT", "STDERR"],
    Remove: true
  },
  Image: "debian:jessie"
  Cmd: ["echo", "hello world!"]
}
```

### Removing a container

```javascript
{
  Dockercmd: "remove",
  Options: {
    Container: "hello_world"
  }
}
```