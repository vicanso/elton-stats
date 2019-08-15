# elton-stats

[![Build Status](https://img.shields.io/travis/vicanso/elton-stats.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-stats)

Route handle stats middleware for elton, it can get some information of route handle, such as status, consuming, size and etc.

```go
package main

import (
	"fmt"
	"net/url"

	"github.com/vicanso/elton"

	proxy "github.com/vicanso/elton-proxy"
	stats "github.com/vicanso/elton-stats"
)

func main() {
	d := elton.New()

	target, _ := url.Parse("https://www.baidu.com")

	d.Use(stats.New(stats.Config{
		OnStats: func(info *stats.Info, _ *elton.Context) {
			fmt.Println(info)
		},
	}))

	d.GET("/*url", proxy.New(proxy.Config{
		Target: target,
		Host:   "www.baidu.com",
	}))

	d.ListenAndServe(":7001")
}
```