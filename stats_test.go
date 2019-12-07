package stats

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vicanso/hes"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestNoStatsPanic(t *testing.T) {
	assert := assert.New(t)
	done := false
	defer func() {
		r := recover()
		assert.Equal(r.(error), errNoStatsFunction)
		done = true
	}()
	New(Config{})
	assert.True(done)
}

func TestSkip(t *testing.T) {
	assert := assert.New(t)
	fn := New(Config{
		OnStats: func(info *Info, _ *elton.Context) {

		},
	})
	c := elton.NewContext(nil, nil)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	c.Committed = true
	err := fn(c)
	assert.Nil(err)
	assert.True(done)
}

func TestStats(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "http://127.0.0.1/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		c.BodyBuffer = bytes.NewBufferString("abcd")
		done := false
		fn := New(Config{
			OnStats: func(info *Info, _ *elton.Context) {
				if info.Status != http.StatusOK {
					t.Fatalf("status code should be 200")
				}
				done = true
			},
		})
		c.Next = func() error {
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.True(done)
	})

	t.Run("return hes error", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "http://127.0.0.1/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		done := false
		fn := New(Config{
			OnStats: func(info *Info, _ *elton.Context) {
				assert.Equal(info.Status, http.StatusBadRequest)
				done = true
			},
		})
		c.Next = func() error {
			return hes.New("abc")
		}
		err := fn(c)
		assert.NotNil(err)
		assert.True(done, "on stats shouldn be called when return error")
	})

	t.Run("return normal error", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "http://127.0.0.1/users/me", nil)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		done := false
		fn := New(Config{
			OnStats: func(info *Info, _ *elton.Context) {
				assert.Equal(info.Status, http.StatusInternalServerError)
				done = true
			},
		})
		c.Next = func() error {
			return errors.New("abc")
		}
		err := fn(c)
		assert.NotNil(err)
		assert.True(done, "on stats shouldn be called when return error")
	})
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
