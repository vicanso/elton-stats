// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stats

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

var (
	errNoStatsFunction = errors.New("require on stats function")
)

type (
	// OnStats on stats function
	OnStats func(*Info, *elton.Context)
	// Config stats config
	Config struct {
		OnStats OnStats
		Skipper elton.Skipper
	}
	// Info stats's info
	Info struct {
		CID        string        `json:"cid,omitempty"`
		IP         string        `json:"ip,omitempty"`
		Method     string        `json:"method,omitempty"`
		Route      string        `json:"route,omitempty"`
		URI        string        `json:"uri,omitempty"`
		Status     int           `json:"status,omitempty"`
		Consuming  time.Duration `json:"consuming,omitempty"`
		Type       int           `json:"type,omitempty"`
		Size       int           `json:"size,omitempty"`
		Connecting uint32        `json:"connecting,omitempty"`
	}
)

// New create a new stats middleware
func New(config Config) elton.Handler {
	if config.OnStats == nil {
		panic(errNoStatsFunction)
	}
	var connectingCount uint32
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	return func(c *elton.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}
		atomic.AddUint32(&connectingCount, 1)
		defer atomic.AddUint32(&connectingCount, ^uint32(0))

		startedAt := time.Now()

		req := c.Request
		uri, _ := url.QueryUnescape(req.RequestURI)
		if uri == "" {
			uri = req.RequestURI
		}
		info := &Info{
			CID:        c.ID,
			Method:     req.Method,
			Route:      c.Route,
			URI:        uri,
			Connecting: connectingCount,
			IP:         c.RealIP(),
		}

		err = c.Next()

		info.Consuming = time.Since(startedAt)
		status := c.StatusCode
		if err != nil {
			he, ok := err.(*hes.Error)
			if ok {
				status = he.StatusCode
			} else {
				status = http.StatusInternalServerError
			}
		}
		if status == 0 {
			status = http.StatusOK
		}
		info.Status = status
		info.Type = status / 100
		size := 0
		if c.BodyBuffer != nil {
			size = c.BodyBuffer.Len()
		}
		info.Size = size

		config.OnStats(info, c)
		return
	}
}
