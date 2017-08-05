# Part 6: Deploy the project on a container service

## What needs to be done.
 - timeout and closed request
 - deploy the container on Heroku

## Code

### Timeout and Closed request.

For timeout and closed request, we use the package context and goroutines.

```go
	ctxTimeout, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
```

In a select statement, if the timeout is reached or the request is closed. We finish it.

```go
	for {
		select {
		case <-ctxTimeout.Done():
            // timeout, closed
            ...
		case c := <-rc:
            // our goroutine
            ...
		}
	}
```

### Deploy the project to heroku

In the Makefile, we add 2 new commands.  
One is for the `multistage build` and the other is for the push on heroku.

We use the multistage build, to have a tiny binary ( from 300mo to 30mo )
Much easier to deploy.

```go
prod:
	docker build -t 6-image-server:latest -f dockerfile.prod.yml .
	docker run \
	-p 8080:8080 \
	--env-file .env 6-image-server:latest

push: build
	heroku container:login
	heroku container:push web --app app-name
```

All the `secrets` will be environnent variables.

```
    // The port is needed by heroku
    PORT=8080
    JWT_SECRET="big-secret-token"
    AWS_BUCKET="bucket-name"
```

