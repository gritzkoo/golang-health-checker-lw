# Migration guidelines

To migrate from: `golang-health-checker` to: `golang-health-checker-lw.v1` you'll need to delacre one functions for each of your integration list.

Because the main package has installed `REDIS` and `MEMCACHED` I'll assume that you project has your own version of this dependencies already installed, so let's get into it!

## Defining a Redis handle function

In your app/pkg where you set your redis context declare a function like below:

```go
package yourredis

import (

    "github.com/go-redis/redis" // assuming you are using this package

    "github.com/gritzkoo/golang-health-checker-lw/pkg/healthchecker"
)

func RedisTest() healthchecker.CheckResponse {
    result := healthchecker.CheckResponse{
        URL: "host-to-my-redis-instance:port"
    }
    rdb := redis.NewClient(&redis.Options{
        Addr:     result.URL,
        // and any other config you need/whant
    })
    response, err := rdb.Ping().Result() // it's important to test only connectivity
    rdb.Close()
    if err != nil {
        result.Error = err
    }
    if response != "PONG" {
        result.Error = fmt.Errorf("Redis client return a nom PONG anwer! got: %s", response)
    }
    return result
}
```

## Defining a memcache handle function

In your app/pkg where you set your memcache context declare a function like below:

```go
package yourmemcache

import (
    "github.com/bradfitz/gomemcache/memcache" // assuming you are using this package

    "github.com/gritzkoo/golang-health-checker-lw/pkg/healthchecker"

)

func MemcacheTest() healthchecker.CheckResponse{
    result := healthchecker.CheckResponse{
        URL: "your-memcache-host:port"
    }
    mcClient := memcache.New(result.URL)
    result.Error = mcClient.Ping()
    return result
}
```

## Then create a healthchecker actor and replace in you http interface

``` go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "yourmodname/yourredis" // here goes the name where you declared the Redis Handle function
    "yourmodname/yourmemcache" // here goes the name where you declared the memcached Handle function

    "github.com/gritzkoo/golang-health-checker-lw/pkg/healthchecker"
)

var checker = healthchecker.New(healthchecker.Config{
    Name:    "My app name", // [optional]
    Version: "v1.0.0", // or get it from some where and replace it here! [optional]
    Integrations: []healthchecker.Check{
        {
            Name:   "My redis integration",
            Handle: yourredis.RedisTest, // invoke your function
        },
        {
            Name:   "My memcache integration",
            Handle: yourmemcache.MemcacheTest, // invoke your function
        },
    },
})
func main() {
    http.HandleFunc("/health-check/liveness", func(w http.ResponseWriter, r *http.Request) {
        resp, _ := json.MarshalIndent(checker.Liveness(), "", "  ")
        w.Header().Add("Content-type", "application/json")
        w.Write(resp)
    })
    http.HandleFunc("/health-check/readiness", func(w http.ResponseWriter, r *http.Request) {
        check := checker.Readiness()
        if !check.Status {
            // do something like write some log or call
            // other service passing full information
            // to handle this issue
        }
        resp, _ := json.MarshalIndent(check, "", "  ")
        w.Header().Add("Content-type", "application/json")
        w.Write(resp)
    })
    http.ListenAndServe(":8090", nil)
}
```
