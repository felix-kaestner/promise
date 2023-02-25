# Promise

<p align="center">
    <span>Go Promise Implementation with support for Generics (requires Go v1.18+).</span>
    <br>
    <span>Run async operations lazily in a separate goroutine on the fly.</span>
    <br><br>
    <a href="https://github.com/felix-kaestner/promise/issues">
        <img alt="Issues" src="https://img.shields.io/github/issues/felix-kaestner/promise?color=29b6f6&style=flat-square">
    </a>
    <a href="https://github.com/felix-kaestner/promise/stargazers">
        <img alt="Stars" src="https://img.shields.io/github/stars/felix-kaestner/promise?color=29b6f6&style=flat-square">
    </a>
    <a href="https://github.com/felix-kaestner/promise/blob/main/LICENSE">
        <img alt="License" src="https://img.shields.io/github/license/felix-kaestner/promise?color=29b6f6&style=flat-square">
    </a>
    <a href="https://pkg.go.dev/github.com/felix-kaestner/promise">
        <img alt="Stars" src="https://img.shields.io/badge/go-documentation-blue?color=29b6f6&style=flat-square">
    </a>
    <a href="https://goreportcard.com/report/github.com/felix-kaestner/promise">
        <img alt="Issues" src="https://goreportcard.com/badge/github.com/felix-kaestner/promise?style=flat-square">
    </a>
    <a href="https://codecov.io/gh/felix-kaestner/promise">
        <img src="https://img.shields.io/codecov/c/github/felix-kaestner/promise?style=flat-square&token=fkA8YwGXkk"/>
    </a>
    <a href="https://twitter.com/kaestner_felix">
        <img alt="Twitter" src="https://img.shields.io/badge/twitter-@kaestner_felix-29b6f6?style=flat-square">
    </a>
</p>

## Features

* Easy interface for composing async operations
* Executes a function in a separate goroutine
* Error Handling using functions (also executed in a separate goroutine)
* Support for Generics (requires Go v1.18+)
* Promises are resolved **lazily**, upon a first call to `Await`, `AwaitOr`, `Then`, `OnSuccess` or `onFailure`
* Support for `promise.All` and `promise.Race` (equivalent to the JavaScript [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) object) 

## Quickstart

```go
package main

import (
	"log"
	"net/http"

	"github.com/felix-kaestner/promise"
)

func main() {
    // Create a new promise.
    // In this example the http request is executed in a separate goroutine
    p := promise.New(func() (*http.Response, error) {
        return http.Get("https://jsonplaceholder.typicode.com/posts/1")
    })
    
    // Handle successful and failed operations in a separate goroutine
    p.Then(func(res *http.Response) {
        log.Printf("Status: %s", res.Status)
    }, func(err error) {
        log.Fatalln(err)
    })

    // Handle only successful operations in a separate goroutine
    p.onSuccess(func(res *http.Response) {
        log.Printf("Status: %s", res.Status)
    })

    // Handle only failed operations in a separate goroutine
    p.onFailure(func(err error) {
        log.Fatalln(err)
    })

    // Await the promise.
    // This blocks execution until the promise is resolved.
    res, err := p.Await()
    
    // Provide a default value (calls Await() internally).
    res = p.AwaitOr(nil)

    // Use channels to select the awaited promise 
    select {
    case <-p.Done():
        res, err = p.Await() // returns immediately since the promise is already resolved
    case <-time.After(5000 * time.Millisecond):
        fmt.Println("Timeout")
    }

    // Take multiple promises and wait for all of them to be finished
    p1 := promise.New(func() (*http.Response, error) {
        return http.Get("https://jsonplaceholder.typicode.com/posts/1")
    })
    p2 := promise.New(func() (*http.Response, error) {
        return http.Get("https://jsonplaceholder.typicode.com/posts/2")
    })
    res, err := promise.All(p1, p2).Await()

    // Take multiple promises and wait until the first of them to is finished
    p1 := promise.New(func() (*http.Response, error) {
        return http.Get("https://jsonplaceholder.typicode.com/posts/3")
    })
    p2 := promise.New(func() (*http.Response, error) {
        return http.Get("https://jsonplaceholder.typicode.com/posts/4")
    })
    res, err := promise.Race(p1, p2).Await()

}
```

##  Installation

Install with the `go get` command:

```
$ go get -u github.com/felix-kaestner/promise
```

## Contribute

All contributions in any form are welcome! 🙌🏻  
Just use the [Issue](.github/ISSUE_TEMPLATE) and [Pull Request](.github/PULL_REQUEST_TEMPLATE) templates and I'll be happy to review your suggestions. 👍

---

Released under the [MIT License](LICENSE).
