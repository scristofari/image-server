build:
	docker build -t 6-image-server .

dev: build
	docker run -it \
	-v $(PWD):/go/src/github.com/scristofari/image-server/part6 \
	-p 8080:8080 \
	--env-file .env 6-image-server:latest bash

prod:
	docker build -t 6-image-server:latest -f dockerfile.prod.yml .
	docker run \
	-p 8080:8080 \
	--env-file .env 6-image-server:latest

push: build
	heroku container:login
	heroku container:push web --app app-name
