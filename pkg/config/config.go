// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

type ExternalMetricsConfig struct {
	Rules []MetricRule `yaml:"rules"`
}

// MetricRule describes rules for transforming Wavefront metrics to/from external metrics API resources.
type MetricRule struct {

	// Query specifies a Wavefront ts query
	Query string `yaml:"query"`

	// The unique name to assign to this metric rule
	Name string `yaml:"name"`
}
