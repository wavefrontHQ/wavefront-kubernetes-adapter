package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type ClientMock struct {
	called  bool
	lastReq http.Request
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	c.called = true
	c.lastReq = *req
	return &http.Response{}, nil
}

func TestDefaultWavefrontClient_Do(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		baseUrl, _ := url.Parse("https://base.url")
		clientMock := &ClientMock{}
		wfClient := NewWavefrontClient(baseUrl, "of good news", 7*time.Second)

		clientRef := wfClient.(*DefaultWavefrontClient)
		clientRef.client = clientMock
		wfClient.Do("GET", "foo", url.Values{
			"l": {"500"},
		})

		assert.True(t, clientMock.called)
		assert.Equal(t, clientMock.lastReq.Method, "GET")
		assert.Equal(t, clientMock.lastReq.URL.String(), "https://base.url/foo?l=500")
		assert.Equal(t, clientMock.lastReq.Header.Values("Authorization"),
			[]string{
				"Bearer of good news",
			})
	})

	t.Run("TODO: Unhappy path", func(t *testing.T) {
		baseUrl, _ := url.Parse("https://base.url")
		clientMock := &ClientMock{}
		wfClient := NewWavefrontClient(baseUrl, "of good news", 7*time.Second)

		clientRef := wfClient.(*DefaultWavefrontClient)
		clientRef.client = clientMock
		wfClient.Do("GET", "foo", url.Values{
			"l": {"500"},
		})

		assert.True(t, clientMock.called)
		assert.Equal(t, clientMock.lastReq.Method, "GET")
		assert.Equal(t, clientMock.lastReq.URL.String(), "https://base.url/foo?l=500")
		assert.Equal(t, clientMock.lastReq.Header.Values("Authorization"),
			[]string{
				"Bearer of good news",
			})
	})
}

func TestNewWavefrontClient(t *testing.T) {
	baseUrl, _ := url.Parse("https://base.url")

	type args struct {
		baseURL *url.URL
		token   string
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want WavefrontClient
	}{
		{
			name: "reasonable timeout",
			args: args{
				baseURL: baseUrl,
				timeout: 3 * time.Second,
				token:   "whatever",
			},
			want: &DefaultWavefrontClient{
				baseURL: baseUrl,
				token:   "whatever",
				client:  &http.Client{Timeout: 3 * time.Second},
			},
		},
		{
			name: "zero timeout yields client with default timeout",
			args: args{
				baseURL: baseUrl,
				timeout: 0 * time.Second,
				token:   "whatever",
			},
			want: &DefaultWavefrontClient{
				baseURL: baseUrl,
				token:   "whatever",
				client:  &http.Client{Timeout: 10 * time.Second},
			},
		},
		{
			name: "negative timeout yields client with default timeout",
			args: args{
				baseURL: baseUrl,
				timeout: -4 * time.Second,
				token:   "whatever",
			},
			want: &DefaultWavefrontClient{
				baseURL: baseUrl,
				token:   "whatever",
				client:  &http.Client{Timeout: 10 * time.Second},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWavefrontClient(tt.args.baseURL, tt.args.token, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWavefrontClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
