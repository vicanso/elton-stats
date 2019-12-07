package main

import (
	"bytes"
	"encoding/json"
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
