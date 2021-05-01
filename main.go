package main

import (
	"context"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	adapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var App *fiber.App
var fiberLambda *adapter.FiberLambda

func init() {
	godotenv.Load()
	App = fiber.New()
	App.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma format e.g. "localhost, nikschaefer.tech"
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	database.ConnectDB()

	router.Initalize(App)

	fiberLambda = adapter.New(App)
}
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}

// func run(){
// 	err := App.Listen(":3000")
// 	if err != nil {
// 		lambda.Start(Handler)
// 	}
// }
func main() {
	// run()
	lambda.Start(Handler)
}

/*
Set Env variables, PORT will auto set to 3000 if not set
*
PORT=3000
DATABASE_URL="host=localhost port=5432 user=postgres password= dbname= sslmode=disable"
*/

/*
Docker Command for Postgres database
*
docker run --name database -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:alpine
*
DB_URL Variable for docker database
*
DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
*/

/*
Generate Zip file for AWS Lambda

$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build -o main
~\Go\Bin\build-lambda-zip.exe -output main.zip main

will generate a main.zip file to upload to lambad
*/
