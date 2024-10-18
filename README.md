# URLShortenerService

URLShortenerService is an API that will turned a long URL into a short one.

## Build the application

You can build the application using the command: 

```
docker build -t url-shortener-service .
```

## Run the application

You can run the application using the command:

```
docker-compose up --build
```

After the application has been shutdown, you can run the command:
```
docker-compose down
```

## Check the health of the application

In order to ensure that the application is up and running, you can do the following cURL:

```
curl -i http://localhost:8080/health
```

The expected status code is 200 OK

## Run the test

In order to run the unitary test, you should use docker-compose then run the test.

```
docker-compose -f docker-compose.test.yml up -d
go test ./...
docker-compose -f docker-compose.test.yml down
```