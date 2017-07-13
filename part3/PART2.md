## Part 2: A minimalist web server to upload and get images.

## Urls:

    - a upload url
    - a url to get the image.

## Dependency management tool

As dependency management tool, we will use the new
golang tool "dep" which will ensure all the dependencies 
of the project.

This will be necessary to get the dependencies for this part.
In this case, the router "gorilla/mux".

## What's needed to be done.
    - a web server
    - a upload url
    - a get image url
    - each image must be in a images folder
    - each image must be less than 5Mo
    - eaxh image must be of the minetype jpg

## Command
    - dep init