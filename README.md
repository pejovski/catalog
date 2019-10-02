# Catalog

## Requirements

```bash
go version 1.13
docker
docker-compose
golang-statik: sudo apt install golang-statik
```

## Installation

Use [go mod](https://blog.golang.org/using-go-modules) to install dependencies.

```bash
go mod tidy
```

Run [docker-compose](https://docs.docker.com/compose/) to build docker images and run necessary containers.

```bash
docker-compose up -d
```

if the command above returns this error

```bash
es01 exited with code 78
```

that is because

```bash
elasticsearch     | [1]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
```

Therefore we need to increase the vm.max_map_count limit:

```bash
sudo sysctl -w vm.max_map_count=524288
```

Now we need to edit /etc/sysctl.conf so the setting will also be in effect after a reboot.

Look for any vm.max_map_count line in /etc/sysctl.conf. If you find one, set its value to 524288. If there is no such line present, add the line

```bash
vm.max_map_count=524288
```
to the end of /etc/sysctl.conf

### Dependency 
- Event Bus - RabbitMQ [Common](https://github.com/pejovski/common)
- Wish List API - [Wish List](https://github.com/pejovski/wish-list)

### Usage
- Make sure the shared RabbitMQ container is up and running [Common](https://github.com/pejovski/common)
- Make sure [Wish List API](http://localhost:8203) is active from [Wish List](https://github.com/pejovski/wish-list)
```bash
docker-compose up -d
go run main.go
```
- open [Catalog API](http://localhost:8201)
- play!
- add new product, add wish list item, update price, update product, etc.

## Swagger update
- use http://editor.swagger.io
- modify app/swagger/swagger.yaml
- run: statik -src=./app/swagger -dest=./app

# Architecture and Design

The project code follows the design principles from the resources bellow

### Microsoft Micro-Services

https://docs.microsoft.com/en-us/dotnet/architecture/microservices/index

### Uber Go Code Structure

https://www.youtube.com/watch?v=nLskCRJOdxM
- extended with receivers and emitters for working with events

### Rest HTTP Server by Go veteran

https://www.youtube.com/watch?v=rWBSMsLG8po

## License
[MIT](https://choosealicense.com/licenses/mit/)