// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/config"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

type RuleHandlerFunc func([]config.MetricRule)

type ExternalMetricsDriver interface {
	getMetricNames() []string
	getQuery(metric string) string
	registerListener(listener ExternalConfigListener)
}

type WavefrontExternalDriver struct {
	cfgFile    string
	rules      map[string]config.MetricRule
	lock       sync.RWMutex
	cfgModTime time.Time
	listener   ExternalConfigListener
}

func NewExternalMetricsDriver(client kubernetes.Interface, cfgFile string) ExternalMetricsDriver {
	driver := &WavefrontExternalDriver{
		cfgFile: cfgFile,
		rules:   make(map[string]config.MetricRule),
	}
	StartHPAListener(client, driver.addRules, driver.deleteRules)
	if cfgFile != "" {
		driver.loadConfig()
	}
	return driver
}

func (d *WavefrontExternalDriver) loadConfig() {
	go wait.Until(func() {
		fileInfo, err := os.Stat(d.cfgFile)
		if err != nil {
			log.Fatalf("unable to get external config file stats: %v", err)
		}

		if fileInfo.ModTime().After(d.cfgModTime) {
			metricsConfig, err := config.FromFile(d.cfgFile)
			if err != nil {
				log.Fatalf("unable to load external metrics discovery configuration: %v", err)
			}
			d.cfgModTime = fileInfo.ModTime()
			d.addRules(metricsConfig.Rules)
		}
	}, 1*time.Minute, wait.NeverStop)
}

func (d *WavefrontExternalDriver) addRules(rules []config.MetricRule) {
	if len(rules) == 0 {
		return
	}

	d.lock.Lock()
	for _, rule := range rules {
		d.rules[rule.Name] = rule
	}
	d.lock.Unlock()

	// always release lock before notifying listeners
	if d.listener != nil {
		d.listener.configChanged()
	}
	log.Debugf("added external metrics rules: %v", rules)
}

func (d *WavefrontExternalDriver) deleteRules(rules []config.MetricRule) {
	if len(rules) == 0 {
		return
	}

	d.lock.Lock()
	for _, rule := range rules {
		delete(d.rules, rule.Name)
	}
	d.lock.Unlock()

	// always release lock before notifying listeners
	if d.listener != nil {
		d.listener.configChanged()
	}
	log.Debugf("deleted external metrics rules: %v", rules)
}

func (d *WavefrontExternalDriver) registerListener(listener ExternalConfigListener) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	d.listener = listener
	log.Info("external configuration listener registered")
}

func (d *WavefrontExternalDriver) getMetricNames() []string {
	d.lock.RLock()
	defer d.lock.RUnlock()

	keys := make([]string, 0, len(d.rules))
	for k := range d.rules {
		keys = append(keys, k)
	}
	return keys
}

func (d *WavefrontExternalDriver) getQuery(metric string) string {
	d.lock.RLock()
	defer d.lock.RUnlock()

	query, found := d.rules[metric]
	if !found {
		return ""
	}
	return query.Query
}
