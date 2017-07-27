build:
	docker build -t 5-image-server .

dev: build
	docker run -it \
	-v $(PWD):/go/src/github.com/scristofari/image-server/part5 \
	-p 8080:8080 \
	5-image-server:latest bash
