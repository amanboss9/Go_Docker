# Go_Docker
Golang production code boilerplate with Redis, MongoDB, Docker and Logger including the best practices.

### Pre Requisite:
- docker
- docker image `mongo`

### Start MongoDB on docker:
`sudo docker run -v "$(pwd)":/data --name mongo -d mongo mongod --smallfiles`

`sudo docker run -it --link mongo:mongo --rm mongo sh -c 'exec mongo "$MONGO_PORT_27017_TCP_ADDR:$MONGO_PORT_27017_TCP_PORT/test"'`


### Make Go build:

- Before creating the build, Check the `DbHost` in conf.json and change it to **mongo** address in docker which can be found with below command:
`sudo docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongo`

`go build`

### Docker Build :

`sudo docker build -t go_docker .`

`sudo docker run --name=go_docker -d -p 9003:9003 go_docker:latest`

### Run and Check:

- Please check your `dockerhostAddress` with the help of below command:

`sudo docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' container-id`

- To get the `container-id` use below docker command:

`sudo docker ps`

`http://dockerhostAddress:9235/one`


### Other important command:

- List all containers(ID)
`docker ps -aq`
- Stop all containers
`docker stop $(docker ps -aq)`
- Stop indivisual container
`docker stop container-id`
- Remove all containers
`docker rm $(docker ps -aq)`
- Remove already in use container
`docker rm container-id`
- Remove all Images
`docker rmi $(docker images -q)`



