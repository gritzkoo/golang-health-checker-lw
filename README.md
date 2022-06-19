# Golang Health Checker lightweight

[![test](https://github.com/gritzkoo/golang-health-checker-lw/actions/workflows/test.yaml/badge.svg)](https://github.com/gritzkoo/golang-health-checker-lw/actions/workflows/test.yaml)
[![build](https://github.com/gritzkoo/golang-health-checker-lw/actions/workflows/block.yaml/badge.svg)](https://github.com/gritzkoo/golang-health-checker-lw/actions/workflows/block.yaml)
[![Coverage Status](https://coveralls.io/repos/github/gritzkoo/golang-health-checker-lw/badge.svg?branch=main)](https://coveralls.io/github/gritzkoo/golang-health-checker-lw?branch=main)
[![CodeQL](https://github.com/gritzkoo/golang-health-checker-lw/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/gritzkoo/golang-health-checker-lw/actions/workflows/codeql-analysis.yml)
![GitHub](https://img.shields.io/github/license/gritzkoo/golang-health-checker-lw)
[![Go Reference](https://pkg.go.dev/badge/github.com/gritzkoo/golang-health-checker-lw.svg)](https://pkg.go.dev/github.com/gritzkoo/golang-health-checker-lw)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gritzkoo/golang-health-checker-lw)](https://img.shields.io/github/go-mod/go-version/gritzkoo/golang-health-checker-lw)
![GitHub repo size](https://img.shields.io/github/repo-size/gritzkoo/golang-health-checker-lw)
[![Go Report Card](https://goreportcard.com/badge/github.com/gritzkoo/golang-health-checker-lw)](https://goreportcard.com/report/github.com/gritzkoo/golang-health-checker-lw)
![GitHub issues](https://img.shields.io/github/issues-raw/gritzkoo/golang-health-checker-lw)
![GitHub Release Date](https://img.shields.io/github/release-date/gritzkoo/golang-health-checker-lw)

___

## Welcome!

This package is a younger brother of the [Golang Health Checker](https://github.com/gritzkoo/golang-health-checker) package that includes only a way to standardize checks of any kind of application and its integrations.

The main purpose of this repository is to disponibility a more simple and lightweight package, with no dependencies, avoiding security issues and extra install packages.

If you are already familiar with `Golang Health Checker`, this package has the same way of working, you need to create a list of tests to create an endpoint in your API to test your integrations, and get an interface with all status for each integration you passes to this package.

___

## How can I install this package?

```sh
go get github.com/gritzkoo/golang-health-checker-lw
```

## How can I use this package?

You can actually test any kind of thing, you just need to declare a function in your application that returns a `healthchecker.CheckResponse` struct and `DONE!`

Below is an `example` of **_HOW TO_** create a `Handle` function to inject into this package.

```go
package main

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gritzkoo/golang-health-checker-lw/pkg/healthchecker"
)
// declaring a function that will test somenting
// this example test a http request calling
// https://github.com/status to check if GitHub is up and running
func TestMyApi() healthchecker.CheckResponse {
  result := healthchecker.CheckResponse{
    URL: "http://github.com/status",
  }
  client := http.Client{}
  request, err := http.NewRequest("GET", result.URL, nil)
  if err != nil {
    result.Error = err
    return result
  }
  response, err := client.Do(request)
  if err != nil {
    result.Error = err
    return result
  }
  if response.StatusCode != http.StatusOK {
    result.Error = fmt.Errorf("The API returned a status code different of 200! code: %d", response.StatusCode)
  }
  return result
}
// create a pointer with for your application of HealthChecker actor
// to be used in your http interface
var checker = healthchecker.New(healthchecker.Config{
  Name:    "My app name", // optional parameter, should be your application name
  Version: "v1.0.0", // or get it from some where and replace it here!
  Integrations: []healthchecker.Check{ // the list of things you need to check
    {
      Name:   "Github integration",
      Handle: TestMyApi,
    },
  },
})

func main() {
  // declaring a endpoint to liveness
  http.HandleFunc("/health-check/liveness", func(w http.ResponseWriter, r *http.Request) {
    resp, _ := json.MarshalIndent(checker.Liveness(), "", "  ")
    w.Header().Add("Content-type", "application/json")
    w.Write(resp)
  })
  // declaring a endpoint to readiness
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

___

## Migration guidelines

You will find all information you need in [THIS DOC](./docs/migration-guidlines.md)
