package provider

import (
	"github.com/golang/glog"
	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"sync"
	"time"
)

type ExternalMetricsDriver interface {
	loadConfig()
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

func (d *WavefrontExternalDriver) registerListener(listener ExternalConfigListener) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	d.listener = listener
	glog.V(5).Info("External configuration listener registered")
}

func (d *WavefrontExternalDriver) loadConfig() {
	go wait.Until(func() {
		fileInfo, err := os.Stat(d.cfgFile)
		if err != nil {
			glog.Fatalf("unable to get external config file stats: %v", err)
		}

		if fileInfo.ModTime().After(d.cfgModTime) {
			metricsConfig, err := config.FromFile(d.cfgFile)
			if err != nil {
				glog.Fatalf("unable to load external metrics discovery configuration: %v", err)
			}
			d.cfgModTime = fileInfo.ModTime()
			d.loadRules(metricsConfig)
		}
	}, 1*time.Minute, wait.NeverStop)
}

func (d *WavefrontExternalDriver) loadRules(externalCfg *config.ExternalMetricsConfig) {
	if externalCfg.Rules == nil {
		return
	}
	d.lock.Lock()
	d.rules = make(map[string]config.MetricRule)
	for _, rule := range externalCfg.Rules {
		d.rules[rule.Name] = rule
	}
	d.lock.Unlock()

	if d.listener != nil {
		d.listener.configChanged()
	}
	glog.V(5).Info("loaded external metrics rules", d.rules)
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
