build:
	docker build -t 4-image-server .

dev: build
	docker run -it \
	-v $(PWD):/go/src/github.com/scristofari/image-server \
	-p 8080:8080 \
	4-image-server:latest bash

test:
	go test -v