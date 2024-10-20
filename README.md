# URLShortenerService

URLShortenerService is an API that will turned a long URL into a short one.

## Build, run and test

### Build the application

You can build the application using the command: 

```
docker build -t url-shortener-service .
```

### Run the application

You can run the application using the command:

```
docker-compose up --build
```

After the application has been shutdown, you can run the command:
```
docker-compose down
```

### Check the health of the application

In order to ensure that the application is up and running, you can do the following cURL:

```
curl -i http://localhost:8080/health
```

The expected status code is 200 OK

### Run the test

In order to run the unitary test, you should use docker-compose then run the test.

```
docker-compose -f docker-compose.test.yml up -d
go test ./...
docker-compose -f docker-compose.test.yml down
```

## Swagger

Once the application is running, you can access the **Swagger UI** interface by clicking [here](http://localhost:8080/swagger/index.html) or visiting the following URL in your browser: `http://localhost:8080/swagger/index.html`

## Explanation of the Shortened Algorithm

The algorithm takes a URL and generates a short identifier called a slug. It does this by following these steps:

1. **Hashing**: The URL is processed using the SHA-1 hashing function, which creates a unique string of characters.
2. **Base62 Encoding**: The hash is then converted into Base62 format, which uses 62 different characters (0-9, A-Z, a-z) to create shorter strings.
3. **Truncation**: The resulting string is shortened to 8 characters to create the final slug.

### Possibilities and Collisions

With 8 characters in Base62, there are about **218 trillion (62^8)** different possible slugs. Because of this huge number, it's very unlikely that two different URLs will produce the same slug. It's also worth noting that slugs and their original URLs are typically not stored forever. This means that even if a collision happens, the chance of it affecting users is almost 0.

### Consistency

One important feature of this algorithm is consistency. This consistency is achieved because the algorithm uses a deterministic process to generate the slug:

1. The SHA-1 hash is always the same for the same input URL, meaning that there are no random variables like timestamps or salts involved in the hashing process.
2. The Base62 encoding of this hash also does not change, as it strictly follows the same conversion method without any random elements.

This means that the same URL will always generate the same slug each time it is processed. Slug generation will not change for a specific URL, as it is not influenced by any external factors like random numbers or changing variables.