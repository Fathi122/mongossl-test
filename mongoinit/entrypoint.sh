#!/bin/bash

nohup gosu mongodb mongod --dbpath=/data/db &
nohup gosu mongodb mongo admin --eval "help" > /dev/null 2>&1
res=$?

while [[ "$res" -ne 0 ]]; do
  echo "Waiting mongodb starting"
  mongo admin --eval "help" > /dev/null 2>&1
  res=$?
  sleep 2
done

# create user
nohup gosu mongodb mongo mongodbssl --eval "db.createUser({ user: 'testuser', pwd: 'testuser', roles: [{ role: 'dbAdmin', db: 'mongodbssl' },{ role: 'readWrite', db: 'mongodbssl' }]});"
nohup gosu mongodb mongo admin --eval "db.shutdownServer();"

# restart mongod with ssl support by loading config file
nohup gosu mongodb mongod --dbpath=/data/db --config /etc/mongod.conf --bind_ip_all --auth

