# URLShortenerService

URLShortenerService is an API that will turned a long URL into a short one.

## Build the application

You can build the application using the command: 

```
$> docker build -t url-shortener-service .
```

## Run the application

You can run the application using the command:

```
$> docker run -p 8080:8080 url-shortener-service
```

## Check the health of the application

In order to ensure that the application is up and running, you can do the following cURL:

```
$> curl -i http://localhost:8080/health
```

The expected status code is 200 OK
