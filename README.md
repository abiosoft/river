river
=====
River is a simple and lightweight REST server.

### Getting Started
```go
rv := river.New()
```

Use middlewares
```go
rv.Use(river.Logger()) 
```

Create endpoints
```go
e := river.NewEndpoint(). 
    Get("/:id", func(c *river.Context){
        id := c.Param("id")
        ... // fetch data with id
        c.Render(200, data)
    }).
    Post("/", func(c *river.Context){
        ... // process c.Body and store in db
        c.Render(201, data)
    })
    ...

e.Use(MyMiddleware) // endpoint specific middleware
```

Handle endpoints
```go
rv.Handle("/user", e) 
```

Run
```go
rv.Run(":8080")
```

Check [example](https://github.com/abiosoft/river/tree/master/example) code for more.

### Approach
* An endpoint is a REST endpoint with handlers for supported methods.
* All endpoints are handled by a River instance.
* Outputs are rendered via a preset or custom Renderer.
* Middlewares and Renderers can be global or specific to an endpoint.

### Request Flow
Basic flow
```
Request -> Middlewares -> Endpoint -> Renderer
```

Full flow
```
                    Request
                       |
                       |  
                     Router
                    /     \                  
                   /       \
                  /         \
              Found      Not Found / Method Not Allowed
                 \          /
                  \        /
                   \      /
              Global Middlewares
                   /      \
                  /        \
 Endpoint Middlewares    Not Found / Method Not Allowed Handler
        |                       |
        |                       |
     Endpoint                Renderer
        |
        |
     Renderer

```

### Endpoint
Create
```go
e := river.NewEndpoint()
```

Handle Requests
```go
e.Get("/", handler).Post(...).Put(...) // method chaining
e.Handle(method, ...) // for custom request methods
```

Any function can be an handler thanks to dependency injection. 
River is also compatible with `http.Handler`. 
```go
func () {...} // valid
func (c *river.Context) {...} // valid
func (c *river.Context, m MyStruct) {...} // valid
func (w http.ResponseWriter, r *http.Request) {...} // valid
```

JSON helper
```go
func (c *river.Context){
    var users []User
    c.DecodeJSONBody(&users)
    ... // process users
}
```

### Middleware
A middleware is any function that takes in the context.
```go
type Middleware func(c *river.Context)
```

River comes with `river.Recovery()` for panic recovery.  

```go
rv.Use(Middleware) // global
e.Use(Middleware)  // endpoint
```

Middleware determines if request should continue. 
```go
func (c *river.Context){
    ... // do something before
    c.Next()
    ... // do something after
}
```

Any `http.Handler` can also be used as a middleware.
```go
rv.UseHandler(handler)
```

### Service Injection
Registering
```go
var m MyStruct
...
rv.Register(m) // global
e.Register(m)  // endpoint
```

This will be passed as parameter to any endpoint handler that has `MyStruct`
as a function parameter.
```go
func handle(c *river.Context, m MyStruct) { ... }
```

Middlewares can also register request scoped service.
```go
func AuthMiddleware(c *river.Context) {
    var session *Session
    ... // retrieve session
    c.Register(session)
}
```

### Renderer
Renderer takes in data from endpoints and renders the data as response.

`context.Render(...)` renders using the configured Renderer. `JSONRenderer` is one of the available renderers. 

Example Renderer, transform response to JSend format before sending as JSON.
```go
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

Setting a Renderer. When an endpoint Renderer is not set, global Renderer is used.
```go
rv.Renderer(MyRenderer) // global
e.Renderer(MyRenderer)  // endpoint
```

### Custom server
River is an `http.Handler`. You can do without `Run()`.
```go
http.ListenAndServe(":8080", rv)
```

### Router
River uses [httprouter](https://github.com/julienschmidt/httprouter) underneath.

### Contributing
* Create an issue to discuss.
* Send in a PR.

### Why the name "River", a "REST" server ? Can you REST on a River ?
Well, yes. You only need to know how to swim or wear a life jacket. 

### License
Apache 2