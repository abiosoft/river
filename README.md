river
=====
River is a simple and lightweight REST server.

### Getting Started
```go
rv := river.New()
```
* Use middlewares
```
rv.Use(river.Logger()) 
```
* Create endpoints
```
e := river.NewEndpoint(). 
    Get("/:id", func(c *river.Context){
        id := c.Param("id")
        ... // fetch data with id
        c.Render(data, 200)
    }).
    Post("/", func(c *river.Context){
        ... // process c.Body and store in db
        c.Render(m, 201)
    })
    ...

e.Use(MyMiddleware) // endpoint specific middleware
```
* Handle endpoints
```
rv.Handle("/user", e) 
```
* Run
```
rv.Run(":8080")
```

### Approach
* An endpoint is a REST endpoint with supported methods.
* All endpoints are handled by a River instance.
* Outputs are rendered via a preset or custom Renderer.
* Middleware and Renderers can be global or endpoint specific.

### Endpoint
Create
```
e := river.NewEndpoint()
```
Request Handler
```
func (c *river.Context)
```
Handle Requests.
```
e.Get(...).Post(...).Put(...) // method chaining
e.Handle(method, ...) // for custom request methods
```

### Middleware
River comes with `river.Logger()` and `river.Recovery()` for logging and panic recovery.  
```
rv.Use(Middleware) // global
e.Use(Middleware)  // endpoint
```
Middleware determines if request should continue. 
```
func (c *river.Context){
    ... // do something before
    c.Next()
    ... // do something after
}
```
Any `http.Handler` can be a middleware.
```
rv.UseHandler(handler)
```
### Renderer
With River, you can change responses without changing a line of code in your endpoints.
Renderer takes in data from endpoints and writes it to the ResponseWriter.

`c.Render(...)` renders using the configured Renderer. `JSONRenderer` is one of the available renderers. 

Creating a Renderer. e.g. transform response to JSend format before sending as JSON.
```
func MyRenderer (c *river.Context, data interface{}) error {
    resp := river.M{"status" : "success", "data" : data}
    if _, ok := data.(error); ok {
        resp["status"] = "error"
        resp["message"] = data
        delete(resp, "data")
    }
    return JSONRenderer(c, resp)
}
```
Setting a Renderer.
```
rv.Renderer(MyRenderer) // global
e.Renderer(MyRenderer)  // endpoint
```
### Custom server.
River is an `http.Handler`. You can add it manually without `Run()`.
```
http.ListenAndServe(":8080", rv)
```
### Router
River uses [httprouter](https://github.com/julienschmidt/httprouter) underneath.

### Contributing
* Create an issue to discuss.
* Send in a PR.

### I wrote a middleware
Thanks, create a PR to add it to this README. 

### Why the name "River", a "REST" server ? Can you REST on a River ?
Well, yes. You only need to know how to swim or wear a life jacket. 

### License
Apache 2