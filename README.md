# Go Template

Go-chi Upper DB

### Prerequisites

Make sure you have golang installed on your system. [Golang Instalation instruction](https://golang.org/doc/install)

### Installing Dependencies

To run this project you will need to install the third party dependencies.

Lists dependencies

```
go list -m all
```

Install dependencies with

```
go mod download
```

### Runing Application Locally

To run this application use the command below

```
go run cmd/api/*.go 
```

### Runing In Docker
To run this application using docker use the following command(s) below

```
docker-compose up -d 
```

### Running Binary File

To run this application use the command below

```
sudo make build
```

And

```
./bin/api
```

## Running the tests

To run the automated tests for this system

```
go test ./...
```

## Deployment

This project is dockerized and deployed to heroku

```
web: cat <<< $GOOGLE_CREDENTIALS > $GOOGLE_APPLICATION_CREDENTIALS && bin/migrate -dir ./db -env production up && GOOGLE_APPLICATION_CREDENTIALS=$GOOGLE_APPLICATION_CREDENTIALS DATABASE_URL=$DATABASE_URL APP_ENV=$APP_ENV PORT=$PORT bin/api
```

## Built With
* [Chi](https://github.com/go-chi/chi) - The web framework used

## Authors
* **Paul Xue**
* **Jesse Okeya**

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.