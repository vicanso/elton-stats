# cod-stats

[![Build Status](https://img.shields.io/travis/vicanso/cod-stats.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-stats)

Route handle stats middleware for cod, it can get some information of route handle, such as status, consuming, size and etc.

```go
package main

import (
	"fmt"
	"net/url"

	"github.com/vicanso/cod"

	proxy "github.com/vicanso/cod-proxy"
	stats "github.com/vicanso/cod-stats"
)

func main() {
	d := cod.New()

	target, _ := url.Parse("https://www.baidu.com")

	d.Use(stats.New(stats.Config{
		OnStats: func(info *stats.Info, _ *cod.Context) {
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