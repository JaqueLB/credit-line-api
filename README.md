# Credit Line API

This is an API that determines a credit line given a client.

The API was made with the [Gin Gonic Framework](https://github.com/gin-gonic/gin), which is a lightweight and performatic option for working with RESTful APIs.

As a database, this API uses a simple local storage, but it can be extended to other databases through the `IStorage` interface.

## Run Locally

To run the API locally, you can just use the following command in your terminal:

```sh
go run main.go
```

The API will be served at `localhost:8080`. If you want, you can customize the port by setting the Environment Variable called `PORT` to the number that is available on your machine.

## API Documentation

The API docs are at `./openapi.yaml`
