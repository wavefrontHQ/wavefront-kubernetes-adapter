// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

const DEFAULT_TIMEOUT = 10 * time.Second

type WavefrontClient interface {
	Do(verb, endpoint string, query url.Values) (*http.Response, error)
	ListMetrics(prefix string) ([]string, error)
	Query(ts int64, query string) (QueryResult, error)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DefaultWavefrontClient struct {
	baseURL *url.URL
	token   string
	client  httpClient
}

func NewWavefrontClient(baseURL *url.URL, token string, apiTimeout time.Duration) WavefrontClient {
	clientTimeout := apiTimeout
	if apiTimeout <= 0 {
		clientTimeout = DEFAULT_TIMEOUT
	}
	return &DefaultWavefrontClient{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: clientTimeout},
	}
}

const (
	authzHeader         = "Authorization"
	bearer              = "Bearer "
	chartEndpoint       = "/api/v2/chart/api"
	metricsListEndpoint = "/chart/metrics/list"
	queryKey            = "q"
	startTime           = "s"
	granularity         = "g"
	outsideSeries       = "i"
)

func (w DefaultWavefrontClient) Do(verb, endpoint string, query url.Values) (*http.Response, error) {
	u := *w.baseURL
	u.Path = path.Join(u.Path, endpoint)
	u.RawQuery = query.Encode()

	log.Printf("DefaultWavefrontClient.Do, query: %s", u.String())

	req, err := http.NewRequest(verb, u.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set(authzHeader, bearer+w.token)

	resp, err := w.client.Do(req)
	if err != nil {
		return resp, err
	}

	// Check all 2xx HTTP codes
	code := resp.StatusCode
	if code/100 != 2 {
		return resp, fmt.Errorf("error status=%s code=%d", resp.Status, code)
	}
	return resp, nil
}

func (w DefaultWavefrontClient) ListMetrics(prefix string) ([]string, error) {
	log.Debugf("DefaultWavefrontClient.ListMetrics")

	vals := url.Values{}
	vals.Set("m", prefix)
	vals.Set("l", "500")

	resp, err := w.Do("GET", metricsListEndpoint, vals)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body
	var result ListResult
	if err = json.NewDecoder(body).Decode(&result); err != nil {
		return nil, &Error{
			Type: ErrBadResponse,
			Msg:  err.Error(),
		}
	}
	log.Trace("DefaultWavefrontClient.ListMetrics", result.Metrics)
	return result.Metrics, nil
}

func (w DefaultWavefrontClient) Query(start int64, query string) (QueryResult, error) {
	log.Debugf("DefaultWavefrontClient.Query: start=%d, query=%s", start, query)
	if query == "" {
		return QueryResult{}, &Error{
			Type: ErrBadData,
			Msg:  "empty query string",
		}
	}

	vals := url.Values{}
	vals.Set(queryKey, query)
	vals.Set(startTime, strconv.FormatInt(start, 10))
	vals.Set(granularity, "m")
	vals.Set(outsideSeries, "false")

	resp, err := w.Do("GET", chartEndpoint, vals)
	if err != nil {
		return QueryResult{}, err
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body
	var result QueryResult
	if err = json.NewDecoder(body).Decode(&result); err != nil {
		return QueryResult{}, &Error{
			Type: ErrBadResponse,
			Msg:  err.Error(),
		}
	}
	log.Trace("DefaultWavefrontClient.Query", result)
	return result, nil
}
