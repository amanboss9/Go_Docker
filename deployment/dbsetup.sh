#!/bin/bash
# Reference https://github.com/frodenas/docker-mongodb/blob/master/scripts/first_run.sh
mongo GO_DOCKER_DB --eval "db.createUser({ user: 'admin', pwd: 'admin', roles: [ { role: 'readWrite', db: 'GO_DOCKER_DB' },'dbAdmin' ] });"
