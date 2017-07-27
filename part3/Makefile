build:
	docker build -t 3-image-server .

dev: build
	docker run -it \
	-v $(PWD):/go/src/github.com/scristofari/image-server/part3 \
	-p 8080:8080 \
	3-image-server:latest bash
	