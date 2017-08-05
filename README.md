# Image Server

In this tutorial, we will see how to implement a http service for 
uploading images, with presets and ready for the cloud.

The purpose is to have an image resizer with these functionalities :
 - A secured route to have a token with an expiration date.
 - Presets to adjust or add effects to images.
 - Ready to deploy on the cloud (with heroku)

Summary:

- [Part 1 : Initiating the project](https://github.com/scristofari/image-server/blob/master/part1/PART1.md)

features: docker, basic router, server

- [Part 2: A minimalist web server to upload and get images](https://github.com/scristofari/image-server/blob/master/part2/PART2.md)

features: dependencies, gorilla mux, file), tests

- [Part 3: Image manipulation and presets](https://github.com/scristofari/image-server/blob/master/part3/PART3.md)

features: resize, image encoding/decoding, table driven tests

- [Part 4: Secure the service with token](https://github.com/scristofari/image-server/blob/master/part4/PART4.md)

features: basic auth, jwt, functionnal tests, middlewares

- [Part 5: Add storage providers](https://github.com/scristofari/image-server/blob/master/part5/PART5.md)

features: interface

- [Part 6: Deploy the project on a container service](https://github.com/scristofari/image-server/blob/master/part6/PART6.md)

features: docker multistage, goroutine, context, timeout
