 # HTTP
 
 ## Server
 ### http.HandleFunc vs http.Handle
 
- `http.HandlerFunc`: Type-Adapter -> take normal function and transform to `Handler`
- `http.Handler`: Interface with `ServeHTTP(...)`

- `http.HandleFunc`: register route with normal function
- `http.Handle`: register route with `Handler`

### Middleware pattern

- chain multiple middlewares/interceptors
- interceptor takes `Handler` and returns one
  - wrap input with HandlerFunc and add custom logic

```go
// interceptor pattern -> able to chain interceptors
func interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("method", req.Method, "path", req.URL.Path)
		next.ServeHTTP(w, req)
	})
}
```

## Client

- `http.Client`: create new client (e.g. `http.Get` use default client)
- `defer res.Body.Close()` after successful call
