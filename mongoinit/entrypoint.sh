#!/bin/bash

nohup gosu mongodb mongod --dbpath=/data/db &
ret=0
echo $res
while [[ "$ret" != 1 ]]; do
  echo "Waiting mongodb starting"
  ret=$(mongo admin --quiet --eval "db.adminCommand('ping').ok")
  sleep 2
done

# create user
nohup gosu mongodb mongo mongodbssl --eval "db.createUser({ user: 'testuser', pwd: 'testuser', roles: [{ role: 'dbAdmin', db: 'mongodbssl' },{ role: 'readWrite', db: 'mongodbssl' }]});"
nohup gosu mongodb mongo admin --eval "db.shutdownServer();"

# restart mongod with ssl support by loading config file
nohup gosu mongodb mongod --dbpath=/data/db --config /etc/mongod.conf --bind_ip_all --auth

 