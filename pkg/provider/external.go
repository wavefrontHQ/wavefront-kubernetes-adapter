package provider

import (
	"github.com/golang/glog"
	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/config"
)

type ExternalMetricsDriver interface {
	loadRules(rules *config.ExternalMetricsConfig)
	getMetricNames() []string
	getQuery(metric string) string
}

type WavefrontExternalDriver struct {
	rules map[string]config.MetricRule
}

func (d *WavefrontExternalDriver) loadRules(rules *config.ExternalMetricsConfig) {
	if rules.Rules == nil {
		return
	}

	d.rules = make(map[string]config.MetricRule)
	for _, rule := range rules.Rules {
		d.rules[rule.Name] = rule
	}
	glog.V(5).Info("loaded external metrics rules", d.rules)
}

func (d *WavefrontExternalDriver) getMetricNames() []string {
	keys := make([]string, 0, len(d.rules))
	for k := range d.rules {
		keys = append(keys, k)
	}
	return keys
}

func (d *WavefrontExternalDriver) getQuery(metric string) string {
	query, found := d.rules[metric]
	if !found {
		return ""
	}
	return query.Query
}
