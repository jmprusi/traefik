package headers

// Middleware tests based on https://github.com/unrolled/secure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/containous/traefik/config"
	"github.com/containous/traefik/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomRequestHeader(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	header := newHeader(emptyHandler, config.Headers{
		CustomRequestHeaders: map[string]string{
			"X-Custom-Request-Header": "test_request",
		},
	})

	res := httptest.NewRecorder()
	req := testhelpers.MustNewRequest(http.MethodGet, "/foo", nil)

	header.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test_request", req.Header.Get("X-Custom-Request-Header"))
}

func TestCustomRequestHeaderEmptyValue(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	header := newHeader(emptyHandler, config.Headers{
		CustomRequestHeaders: map[string]string{
			"X-Custom-Request-Header": "test_request",
		},
	})

	res := httptest.NewRecorder()
	req := testhelpers.MustNewRequest(http.MethodGet, "/foo", nil)

	header.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test_request", req.Header.Get("X-Custom-Request-Header"))

	header = newHeader(emptyHandler, config.Headers{
		CustomRequestHeaders: map[string]string{
			"X-Custom-Request-Header": "",
		},
	})

	header.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "", req.Header.Get("X-Custom-Request-Header"))
}

func TestSecureHeader(t *testing.T) {
	testCases := []struct {
		desc     string
		fromHost string
		expected int
	}{
		{
			desc:     "Should accept the request when given a host that is in the list",
			fromHost: "foo.com",
			expected: http.StatusOK,
		},
		{
			desc:     "Should refuse the request when no host is given",
			fromHost: "",
			expected: http.StatusInternalServerError,
		},
		{
			desc:     "Should refuse the request when no matching host is given",
			fromHost: "boo.com",
			expected: http.StatusInternalServerError,
		},
	}

	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	header, err := New(context.Background(), emptyHandler, config.Headers{
		AllowedHosts: []string{"foo.com", "bar.com"},
	}, "foo")
	require.NoError(t, err)

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			req := testhelpers.MustNewRequest(http.MethodGet, "/foo", nil)
			req.Host = test.fromHost
			header.ServeHTTP(res, req)
			assert.Equal(t, test.expected, res.Code)
		})
	}
}
