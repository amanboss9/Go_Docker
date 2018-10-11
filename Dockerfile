# Golang Image 
FROM golang:latest

#WORKDIR /app
WORKDIR /go/src/Go_Docker

# Install dependencies
RUN go get gopkg.in/mgo.v2
RUN go get github.com/hhkbp2/go-logging
RUN go get github.com/gin-contrib/cors
RUN go get github.com/gin-contrib/sessions
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/sessions/redis
RUN go get github.com/go-redis/redis
# RUN go get github.com/mediocregopher/radix.v2/redis

ENV SRC_DIR=/go/src/Go_Docker

# Add the source code:copy the local package files to the container workspace
ADD . $SRC_DIR

# Build it:
RUN go install Go_Docker

#RUN ["chmod", "+x"]
ENTRYPOINT /go/bin/Go_Docker

# Document that the service listens on port 9002.
EXPOSE 9003
