# elton-stats

The middleware has been archived, please use the middleware of [elton](https://github.com/vicanso/elton).

[![Build Status](https://img.shields.io/travis/vicanso/elton-stats.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-stats)

Route handle stats middleware for elton, it can get some information of route handle, such as status, consuming, size and etc.

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/vicanso/elton"

	stats "github.com/vicanso/elton-stats"
)

func main() {
	e := elton.New()

	e.Use(stats.New(stats.Config{
		OnStats: func(info *stats.Info, _ *elton.Context) {
			buf, _ := json.Marshal(info)
			fmt.Println(string(buf))
		},
	}))

	e.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("abcd")
		return
	})
	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
```
