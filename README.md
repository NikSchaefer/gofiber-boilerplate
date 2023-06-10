# 👋 Go-Fiber Boilerplate
Golang Rest API boilerplate built with GORM, Go-Fiber, and a PostgreSQL database. Running in a docker container with Hot Reload.

## Quickstart 🚀
To quickly get started with the Go-Fiber Boilerplate, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/NikSchaefer/go-fiber
```

2. Install Dependencies

```bash
cd go-fiber-boilerplate
go mod tidy
```

3. Connect the database: Create a `.env` file and put in a connection string

```bash
DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
```

4. Start the project
```bash
go run main.go
```

Alternatively, you can build a Docker image and run the project in a container, as seen below.

# File Structure 📁
The file structure of the project is divided into five main folders and a main.go file.
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

### Database 🗄️
The database folder initializes the database connection and registers the models. You can add new models by registering them in connect.go. The global DB variable is initialized in database.go.

### Handlers 🤝
The handlers folder defines each request for each model and how it interacts with the database. The functions are mapped by the router to the corresponding URL links.

### Middleware 🛡️
The middleware folder contains a file for each middleware function. The security middleware is applied first to everything in router.go and applies general security middleware to the incoming requests. The JSON middleware serializes the incoming request so that it only allows JSON. Finally, the Authentication middleware is applied individually to requests that require the user to be logged in.

### Router 🛣️
The router file maps each incoming request to the corresponding function in handlers. It first applies the middleware and then groups the requests to each model and finally to the individual function.

### Main.go 🚀
The main.go file reads environment variables and applies the CORS middleware. You can change the allowed request sites in the configuration. It then connects to the database by running the function from database/connect.go and finally initializes the app through the router.


# Debug 🐛
The port can be specified with an environment variable but will default to 3000 if not specified.

## Database 🗄️

to run the database on docker use the following command: 

`docker run --name database -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:alpine`

To connect to the database you just started you can set the enviroment variable of 

`DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"`

## Docker 🐳
Docker build base image in first stage for development

`docker build --target build -t base .`

run dev container

`docker run -p 3000:3000 --mount type=bind,source="C:\Users\schaefer\go\src\fiber",target=/go/src/app --name fiber -td base`

rebuild and run package

`docker exec -it web go run main.go`

stop and remove container

`docker stop fiber; docker rm fiber`

## Recommended 🙌
run a postgres databse in docker and use the [fiber command line](https://github.com/gofiber/cli) to hot reload your application. Note: you can hot reload using docker or the fiber command line

# API Documentation 📖
This project provides a batteries included REST API that allows the user to interact with a PostgreSQL database. The available endpoints are listed below:

### User Endpoints

`POST /users/`
Create a new user. The request should include a JSON payload with the following fields:

- `username`: a string containing the user's username.
- `password`: a string containing the user's password.

`DELETE /users/`
Delete the user. This endpoint requires the user to be authenticated.

`PUT /users/`
Change the user's password. This endpoint requires the user to be authenticated. The request should include a JSON payload with the following fields:

- `oldPassword`: a string containing the user's current password.
- `newPassword`: a string containing the user's new password.

`POST /users/me`
Get information about the user. This endpoint requires the user to be authenticated.

`POST /users/login`
Log in the user. The request should include a JSON payload with the following fields:

- `username`: a string containing the user's username.
- `password`: a string containing the user's password.

`DELETE /users/logout`
Log out the user. This endpoint requires the user to be authenticated.

### Product Endpoints

`POST /products/`
Create a new product. This endpoint requires the user to be authenticated. The request should include a JSON payload with the following fields:

- `name`: a string containing the name of the product.
- `description`: a string containing the description of the product.
- `price`: a float64 containing the price of the product.

`POST /products/all`
Get a list of all products.

`DELETE /products/:id`
Delete a product by ID. This endpoint requires the user to be authenticated.

`POST /products/:id`
Get a product by ID.

`PUT /products/:id`
Update a product by ID. This endpoint requires the user to be authenticated. The request should include a JSON payload with the following fields:

- `name`: a string containing the new name of the product.
- `description`: a string containing the new description of the product.
- `price`: a float64 containing the new price of the product.
If you need more information about the request and response of each endpoint, please check the corresponding function in the handlers folder.

# License 📜

[MIT](https://choosealicense.com/licenses/mit/)
