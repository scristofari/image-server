# Part 1 : Initialisation of the project

## Image

We will use the official docker image for golang.
There are all the libraries needed to build the project.

## Makefile

Just to be have more "plug & play" commands. :)
The command "make dev", will build and run the project.

## Source code

All the source code must be in the "GOPATH" and in a folder "src".
Elsewhere, it will be impossible to build the project.
Here, in the golang image, "GOPATH=/go".

## What's needed to be done.

 - created a webserver which print `upload route`

 ## Buuild the project.

There're different ways to achieve that:

 - go build .
 Will build a binarie for the project.

 - go install
 Will build a binarie and place it in a `bin` folder

 - go run main.go
 Will build and run the project.