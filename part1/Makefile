build:
	docker build -t 1-image-server .

dev: build
	docker run -it \
	-v $(PWD):/go/src/github.com/scristofari/image-server/part1 \
	-w /go/src/github.com/scristofari/image-server/part1 \
	-p 8080:8080 \
	1-image-server:latest bash