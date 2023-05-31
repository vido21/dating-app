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


# Structure

1. auth: This directory contains the authentication-related code. It includes functionalities such as user registration and login.

2. common: The common directory holds shared code and utilities that are used throughout the application. It includes type, utils, middleware adn common functionality.

3. config: The config directory consists of database and jwt config

4. database: The database directory handles the management and interaction with the database. It consists seeder to populate initial data, connection to set up database connection and migrations to set up schemes.

5. models: The models directory defines the data structures and schema used in the application. It represents the entities and their relationships in the database.

6. mocks: The mocks directory is used to store mock objects or mock implementations of dependencies for testing purpose.

7. routes: The routes directory defines the routes and their corresponding handlers. It maps the incoming HTTP requests to the appropriate controller functions. The route files define the API endpoints and their associated middleware functions.

8. test: The test directory contains fature testing or integration testing. It consists common and cons that are used throughout the application.

9. main.go: The main.go file is the entry point of the application. It initializes the necessary components, sets up the server, and starts listening for incoming requests. It might import and configure the required libraries and frameworks.
    
10. premium-packages: The premium-packages component in our application consists of multiple elements, including models, a controller, a service, and unit tests for premium-packages scheme / use case.

11. profiles: The profiles component in our application consists of multiple elements, including models, a controller, a service, and unit tests for profiles scheme / use case.

12. purchases: The purchases component in our application consists of multiple elements, including models, a controller, a service, and unit tests for purchases scheme / use case.

13. swipes: The swipes component in our application consists of multiple elements, including models, a controller, a service, and unit tests for swipes scheme / use case.

14. users: The users component in our application consists of multiple elements, including models, a controller, a service, and unit tests for users scheme / use case.
    

## Running Application

Rename `.env.example` to `.env` and place your database credentials and jwt secret key

```
$ mv .env.example .env
$ docker-compose up -d
$ go mod vendor
$ go run main.go
```

## Building

```
$ go build -v
$ ./dating-app
```

## Testing

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

