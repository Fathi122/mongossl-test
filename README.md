# mongossl-test
## generate mongodb certs
```
./genmongodbcert.sh
```
## build and start mongodb with ssl support
```
docker-compose -f docker-compose-setupmongo.yaml up --force-recreate --build
```

## run go test client with ssl support
```
go run main.go
```
## stop mongodb deployment
```
docker-compose -f docker-compose-setupmongo.yaml down --remove-orphans
```