build:
	docker build -t 2-image-server .

dev: build
	docker run -it \
	-v $(PWD):/go/src/github.com/scristofari/image-server \
	-w /go/src/github.com/scristofari/image-server \
	-p 8080:8080 \
	2-image-server:latest bash