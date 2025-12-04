package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eesys/go-zero/core/logx"
	"github.com/stretchr/testify/assert"
)

func TestMetricsInterceptor(t *testing.T) {
	logx.Disable()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.NotNil(t, req)
	assert.Nil(t, err)
	interceptor := MetricsInterceptor("test", nil)
	req, handler := interceptor(req)
	resp, err := http.DefaultClient.Do(req)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	handler(resp, err)
}
