# Part 1 : Initiating the project

## What's need to be done.

 - create a webserver which has a route that prints `upload route`.

## Docker Image

We will use the official docker image for [golang](https://hub.docker.com/_/golang/).  
There are all the libraries needed to build the project.  

## Makefile

Just to get more `plug & play` commands. :)  
The command `make dev`, will build and run the project.  

## Source code

In go, all the source code must be in the `GOPATH` in a folder `src`.  
Elsewhere, it will be impossible to build the project.  
Here, in the golang image, `GOPATH=/go`.  

 ## Build the project.

There're different ways to achieve that:

 - go build .  
 Will build a binary for the project.

 - go install  
 Will build a binary and place it in a `bin` folder

 - go run main.go  
 Will build and run the project.

 For more info, go to the [go-tooling-workshop](https://github.com/campoy/go-tooling-workshop/tree/master/2-building-artifacts).

 ## Code - Basic Router / Run server. 

 In this part, we use build-in router in go.

 ```go
r := http.NewServeMux()
r.HandleFunc("/upload", uploadHandler)
 ```

 For a given path, the router will call a function which have as interface.

 ```go
func handleFunc(w http.ResponseWriter, r *http.Request)
 ```
 
