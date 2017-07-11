## Part 1 : Create a docker environment for dev

## Image

We will use the official docker image for golang.
There are all the libraries needed to build the project.

In the future, we will just have a binarie, so the alpine
alone will be necessary.
See Part 6, for multi-staging docker.

## Makefile

Just to be have more "plug & play" commands. :)
The command "make dev", will build and run the project.

## Source code

All the source code must be in the "GOPATH" and in a folder "src".
Elsewhere, it will be impossible to build the project.
Here, in the golang image, "GOPATH=/go".
