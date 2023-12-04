# gomite

## Getting started


### Installation
```
go get github.com/codimite/gomite
```

### Basic endpoint
```go
func main(){
  gm := gomite.Gomite{Port: "8080", Handler: AppRecoverHandler(GlobalHandler{})}
  gomite.InitTemplates([]string{"templates"})
  
  // Endpoint
  gm.HandleFunc("/getting-started", func(rw http.ResponseWriter, req *http.Request) {
    rw.WriteHeader(200)
    rw.Write([]byte("Health check"))
  })
  
  gm.Start() // Start the server
}

// utils
type GlobalHandler struct {}
func AppRecoverHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				AppLogger{}.Error("[Critical Error Occured]", "AppRecoverHandler", r)
				var message struct {
					Message string `json:"message,omitempty"`
				}
				message.Message = "Something Went Wrong"
				rw.Header().Add("Content-Type", "application/json")
				rw.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(rw).Encode(message)
				return
			}
		}()
		handler.ServeHTTP(rw, r)
	})
}
```

```AppRecoverHandler``` - Handle panics & stop program from exiting
<br>
## Working with Interceptors 

### Setup server
```go
import(
  "github.com/codimite/gomite"
  interceptor "github.com/codimite/gomite/interceptors"
)

func main(){
  headerInterceptor := interceptor.MiddlewareChain{HeaderCheck()}

  gm := gomite.Gomite{Port: servePort, Handler: AppRecoverHandler(GlobalHandler{})}
  gomite.InitTemplates([]string{"templates"})

  gm.Handle("/header-check", headerInterceptor.Handler(HeaderCheckHandler))
}

...
```
### Add the interceptor middleware
```go
func HeaderCheck() interceptor.MiddlewareInterceptor {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		reqHeader := req.Header.Get("X-Test-Header")
		if reqHeader != "X-Test-Header-Value" {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}
		next(rw, req)
	}
}
```
### Add the controller
```go
func HeaderCheckHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte("ok"))
	return
}
```
### Change read timeout to 20 seconds
```go
gm.Start(WithReadTimeout(20 * time.Second))
```