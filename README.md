# mongossl-test
## generate mongodb certs
```
./genmongodbcert.sh
```
## build and start mongodb with ssl support
```
docker-compose -f docker-compose-setupmongo.yaml up --force-recreate --build
```
## start mongoshell
```
  mongo --host localhost --port 27017 --ssl --sslCAFile ./ca.crt --sslPEMKeyFile ./mongodb.pem -u username -p userpassword --authenticationDatabase mongodbssl
```
## run go test client with ssl support
```
go run main.go
```
## stop mongodb deployment
```
docker-compose -f docker-compose-setupmongo.yaml down --remove-orphans
```