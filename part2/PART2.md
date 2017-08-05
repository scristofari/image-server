# Part 2: A minimalist web server to upload and get images.

## What's need to be done.

    - an upload url
    - a get image url
    - each image must be in a images folder

## Dependency management tool

As dependency management tool, we use the new
golang tool [dep](https://github.com/golang/dep) which ensures that all dependencies 
of the project are correct.

This will be necessary to get the dependencies for this part.
In this case, the router `gorilla/mux`.

We just need to tap `dep init`
A `vendor` folder will be created with all the dependencies and 
a `Gopkg.toml` + `Gopkg.lock` for locking the deps versions.

## Code: Dependencies / Read - Write file / Tests

We define 2 routes:

```go
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/images/{img}", imageHandler).Methods("GET")
```

The first one uploading the file, the other getting it.  
[gorilla mux](https://github.com/gorilla/mux) and the standart router use the same interface 
described in PART1.

Note that we use the `defer` statement there:

```go
	defer image.Close()
```
This will execute the function until the surrounding function returns.
I find this, must more readable than executing it at the end of the function.
Note that, we close the image after checking if the `r.FormFile("image")` returns
a error. if there's an error, the image can be nil and this will panic.

At the end of the uploadHandler, we set the 201 http status code
It must be added before writing the response.

```
    w.WriteHeader(http.StatusCreated)
```

## Tests

In the tests, we create a fake request and responseWriter.
With it, we can call our handleFunc function and test the status code of the 
response.

All the tests function must begin by Test... .

```
func TestUploadImage(t *testing.T) {
    ...
    r, err := http.NewRequest("GET", "http://localhost/images/golang.png", nil)
    ... 
    w := httptest.NewRecorder()
    ...
    the handleFunc we want he test
    ...
    if w.Code != http.StatusOK {
		t.Errorf("Expected %d, get %d", http.StatusOK, w.Code)
		t.Error(w.Body)
	}
```
