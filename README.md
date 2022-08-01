# Go-Fiber Boilerplate
Golang Rest API boilerplate built with GORM, Go-Fiber, and a PostgreSQL database. Running in a docker container with Hot Reload.

# File structure
```py
database/
  connect.go
  database.go
handlers/
  auth.go
  product.go
middleware/
  json.go
  auth.go
  security.go
model/
  user.go
  product.go
  session.go
router/
  router.go
main.go
```

## Database

The database folder holds 2 files. The first file `connect.go` initalizes the database connection and migrates the registered models. If you are looking to add new models make sure to register them here for the database. The second file `database.go` initalizes the global DB variable that is referenced in other files.

## Handlers

This folder is the place that holds the functions for each model. Here you will define each request and how it interacts with the database. These functions are used mapped by the router the the URL links.

## Middleware

The middleware folder contains a file for each middleware function. The security middleware is applied first to everything in `router.go` and applies general security middleware to the incoming requests. The JSON middleware serializes the incoming request so that it only allows JSON. This is applied after the hello world in `router.go`. Finally the Authentication middleware is applied indivually to requests that require the user to be logged in. 

## Router
The router file maps each incoming request to the corresponding function in `handlers`. It first applies the middleware and then groups the requests to each model and finally to the indiviual function.

## Main.go

The main.go file functions by reading for enviroment variables and applying the CORS middleware. You can change the allowed request sites in the configuration. It then connects to the database by running the function from `database/connect.go` and finally initalizes the app through the router.

# Debug

the port can be specified with an enviroment variable but will default to 3000 if not specified.

## Database

to run the database on docker use the following command: `docker run --name database -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:alpine`. and to connect to the database you can set the enviroment variable of `DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"`

## Docker
Docker build base image in first stage for development
`docker build --target build -t base .`

run dev container
`docker run -p 3000:3000 --mount type=bind,source="C:\Users\schaefer\go\src\fiber",target=/go/src/app --name fiber -td base`

rebuild and run package
`docker exec -it web go run main.go`

stop and remove container
`docker stop fiber; docker rm fiber`

## Recommended
run a postgres databse in docker and use the [fiber command line](https://github.com/gofiber/cli) to hot reload your application. Note: you can hot reload using docker or the fiber command line

# Dependencies

Install dependencies with go

`go mod tidy`

# License

[MIT](https://choosealicense.com/licenses/mit/)
