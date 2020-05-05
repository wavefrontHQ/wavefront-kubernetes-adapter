// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
)

// ErrorType is the type of the API error.
type ErrorType string

const (
	ErrBadData     ErrorType = "bad_data"
	ErrTimeout               = "timeout"
	ErrCanceled              = "canceled"
	ErrBadResponse           = "bad_response"
)

// Error is an error returned by the API.
type Error struct {
	Type ErrorType
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Msg)
}

type Timeseries struct {
	Label string
	Host  string
	Tags  map[string]string
	Data  [][]float64
}

// QueryResult represents the response returned by the API.
type QueryResult struct {
	Name       string       `json:"name"`
	Query      string       `json:"query"`
	Timeseries []Timeseries `json:"timeseries"`
}

type ListResult struct {
	Metrics []string `json:"metrics"`
	Limit   int      `json:"limit"`
}
