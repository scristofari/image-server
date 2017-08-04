# Image manipulation

## What's needed to be done.

 - Max upload size of 5MB.
 - Unique ID for each uploaded file.
 - Use an external dependencie for rezise the image.
 - Use a preset to specify the size zns the type of the new image.

## Code :

### Max upload size.
The package http includes a max bytes reader.
We just defined our custom max size.

```go
    const (
        mb            = 1 << 20
    )
    var (
        uploadMaxSize = 5 * mb
    )
    r.Body = http.MaxBytesReader(w, r.Body, int64(uploadMaxSize))
```

### Unique ID.

For that, the package [go.uuid](https://github.com/satori/go.uuid).
We need to generate an unique id for each uploaded file. The filename is not sufficent.
Like we only accept `png`, the generated filename will be:

```go
    filename := uuid.NewV4().String() + ".png"
```

### Resize and get the preset.

The package used is [resize][github.com/nfnt/resize].  
The goal here, is to resize the image like the folllowing rules.

 - ?r=300x200, will resize the image with a width of 300 and a height of 200.  
 - ?t=300x200, will create a thumbnail of the image with a width of 300 and a height of 200.  
 - ?r=300x0, with a height of 0, we follow the ratio.

in this exemple, we assume that we have a caching system. So no need to save the resized
image on the filesystem, just render it, and add a max age header.

```go
	w.Header().Set("Cache-Control", "max-age=3600")
	png.Encode(w, i)
```

## Tests

At the point, the tests will not change a lot.  
We use table driven tests, the benefit of that, is that we test a lot of cases at one time.
Very useful for testing the preset.

```go
	cases := []struct {
		in   string
		code int
	}{
		{"http://localhost/images/golang.png?r=50x0", http.StatusOK},
		{"http://localhost/images/golang.png?r=23x0&t=23x34", http.StatusBadRequest},
		{"http://localhost/images/golang.png", http.StatusBadRequest},
		{"http://localhost/images/golang.png?r=efkgjergx2", http.StatusBadRequest},
		{"http://localhost/images/golang.png?r=2500x4000", http.StatusBadRequest},
	}
    ...
```
