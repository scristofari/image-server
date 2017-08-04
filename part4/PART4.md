# Part 4: Secure

## What's needed to be done.
 - Add a new route which will give a token and be secure by a basic auth.
 - The token must have an expired date.

## Code

### Refactoring.

The code in the main.go of part3 is greater, but adding a lot of new features will make it 
more and more complex to read.

We move all the `domain` features in a package named `resizer`. With the command `godoc`, we can have 
the api of the package.

```
    PACKAGE DOCUMENTATION

    package resizer
        import "."


    VARIABLES

    var (
        UploadMaxSize = 5 * mb
    )

    FUNCTIONS

    func CheckCredentials(user string, password string) error

    func Resize(filename string, q *Query) (image.Image, error)

    func Uploadfile(image multipart.File) (string, error)
        Uploadfile : ___

    TYPES

    type Preset struct {
        Width  uint
        Height uint
    }

    type Query struct {
        Type   string
        Preset *Preset
    }

    func GetQueryFromURL(u *url.URL) (*Query, error)
```

As we want to use it for diffrents purpose, on cli, on the server, all the http things
remains in the `main.go`.

### Basic auth.

To add this fonctionality, we can implement middlewares.  
These are some `http.HandlerFunc` which can be chained.

```go
    func authBasicHandlerFunc(f http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            user, pass, ok := r.BasicAuth()
            if !ok {
                http.Error(w, "failed to get auth basic credentials", http.StatusForbidden)
                return
            }

            err := resizer.CheckCredentials(user, pass)
            if err != nil {
                http.Error(w, "failed to sign in: "+err.Error(), http.StatusForbidden)
                return
            }

            f(w, r)
        }
    }

    ...

	r.HandleFunc("/access/token", authBasicHandlerFunc(accessHandleFunc)).Methods("GET")    
```

### The token.

For that, we use the [jwt-go](github.com/dgrijalva/jwt-go).  
This package helps us to generate a token with a expire date of 1 minute.

```go 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
	})
```

For the validation, a middleware willl be added.  
The same mecanic as the basic auth middleware.  

```go
    func jwtHandlerFunc(f http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            /** vars from gorilla mux empty, in test case, we do not execute the router */
            hash := strings.Split(r.URL.Path, "/")

            token, err := jwt.Parse(hash[2], func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
                }

                return []byte("secret"), nil
            })
            if err != nil {
                http.Error(w, "failed to authenticate: "+err.Error(), http.StatusUnauthorized)
                return
            }

            if !token.Valid {
                http.Error(w, "failed to authenticate, token not valid", http.StatusUnauthorized)
                return
            }

            f(w, r)
        }
    }
```

## Tests 

We now change the tests to reflex the real calls, we need the router with all the middlewares.  
For this, we run a test server, and do the request.

```go
	var server = httptest.NewServer(handlers())
	defer server.Close()

    ...

    res, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Error(err)
		return
    }
    ...
```