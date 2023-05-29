# Dating App

This project ships following features as default
* ORM Integration
* Easy Database Migration
* Authentication Using Jwt
* Feature Testing
* Unit Testing
* Isolate Feature Testing Environment
* Easy dotenv Management
* Easy to Mock all Interfaces
* CORS Configuration

 Inspired from [GO Echo boilerplate](https://github.com/triaton/go-echo-boilerplate).

## Running Your Application

Rename .env.example to .env and place your database credentials and jwt secret key

```
$ mv .env.example .env
$ docker-compose up -d
$ go mod vendor
$ go run main.go
```

## Building Your Application

```
$ go build -v
$ ./dating-app
```

## Testing Your Application

```
$ docker-compose up -d
$ go test -v ./...
```

### Generate Code Coverage

```
$ docker-compose up -d
$ chmod +x generate-test-coverage.sh
$ ./generate-test-coverage.sh
```
This will generate `cover.html` with detailed coverage result.

### Make mock interfaces with mockery
```
$ mockery --all
```

## Import Postman Collection (API's)

Download [Postman](https://www.getpostman.com/) -> Import -> Import From `datingapp.postman_collection`

Postman collection already have testing script

