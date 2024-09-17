# golang-health-checker-lw

## `[1.1.0] - 2024-09-16`

**_Changes_**

- `healthchecker.Integration.Error` struct changed type from `error` to `string`
- `healthchecker.Config.Concurrence` used to control max paralel routines, default 10
- Update module golang version from 1.19 -> 1.23

**_Examples_**

```go
package main

import "github.com/gritzkoo/golang-health-checker-lw"

var checker =  healthchecker.New(healthchecker.healthchecker.Config{
  Name: "example-app",
  Version: "v1.1.0",
  Concurrence: 100, // new attribute to control max paralel routines, default is 10
  Integrations: []healthchecker.Check{
    {
      Name: "example test with error",
      Handler: func() healthchecker.CheckerResponse {
        // do your on call or test modules
        return healthchecker.CheckerResponse{
          URL: "use to display a reference on response"
          Error: fmt.Errorf("message if error")
        }
      }
    },
    {
      Name: "example test with no error",
      Handler: func() healthchecker.CheckerResponse {
        // do your on call or test modules
        return healthchecker.CheckerResponse{
          URL: "use to display a reference on response"
        }
      }
    },
  }
})

func main(){
  fmt.Println(checker.Readiness())
}
```

---

**_output_**

```log
{
  "name": "example-app",
  "status": false,
  "version": "v1.1.0",
  "date": "2024-09-17T11:57:13+02:00",
  "duration": 1.243803958,
  "integrations": [
    {
      "name": "example test with error",
      "status": false,
      "response_time": 1.242491667,
      "url": "use to display a reference on response",
      "error": "message if error"
    },
    {
      "name": "example test with no error",
      "status": true,
      "response_time": 1.242491667,
      "url": "use to display a reference on response"
    },
  ]
}
```

---
