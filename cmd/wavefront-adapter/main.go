/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"net/url"
	"os"
	"runtime"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/util/logs"

	"github.com/golang/glog"
	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	customprovider "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"

	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/client"
	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/provider"
)

var (
	version string
	commit  string
)

type WavefrontAdapter struct {
	basecmd.AdapterBase

	// Message is printed on successful startup
	Message string
	// MetricsRelistInterval is the interval at which list of metrics are fetched from Wavefront
	MetricsRelistInterval time.Duration
	// Wavefront Server URL of the form https://INSTANCE.wavefront.com
	WavefrontServerURL string
	// Wavefront API token with permissions to query points
	WavefrontAPIToken string
	// The prefix for custom kubernetes metrics in Wavefront
	CustomMetricPrefix string
	// The file containing the metrics discovery configuration
	AdapterConfigFile string
}

func (a *WavefrontAdapter) makeProviderOrDie() customprovider.MetricsProvider {
	dynClient, err := a.DynamicClient()
	if err != nil {
		glog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		glog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	waveURL, err := url.Parse(a.WavefrontServerURL)
	if err != nil {
		glog.Fatalf("unable to parse wavefront url: %v", err)
	}
	waveClient := client.NewWavefrontClient(waveURL, a.WavefrontAPIToken)

	metricsProvider, runnable := provider.NewWavefrontProvider(dynClient, mapper, waveClient, a.CustomMetricPrefix, a.MetricsRelistInterval, a.AdapterConfigFile)
	runnable.RunUntil(wait.NeverStop)
	return metricsProvider
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	cmd := &WavefrontAdapter{
		CustomMetricPrefix:    "kubernetes",
		MetricsRelistInterval: 10 * time.Minute,
	}
	cmd.Name = "wavefront-custom-metrics-adapter"
	flags := cmd.Flags()
	flags.DurationVar(&cmd.MetricsRelistInterval, "metrics-relist-interval", cmd.MetricsRelistInterval, ""+
		"interval at which to fetch the list of all available metrics from Wavefront")
	flags.StringVar(&cmd.WavefrontServerURL, "wavefront-url", "",
		"Wavefront url of the form https://INSTANCE.wavefront.com")
	flags.StringVar(&cmd.WavefrontAPIToken, "wavefront-token", "",
		"Wavefront API token with permissions to query for points")
	flags.StringVar(&cmd.CustomMetricPrefix, "wavefront-metric-prefix", cmd.CustomMetricPrefix,
		"Wavefront Kubernetes Metrics Prefix")
	flags.StringVar(&cmd.AdapterConfigFile, "external-metrics-config", "",
		"Configuration file for driving external metrics API")
	flags.StringVar(&cmd.Message, "msg", "starting wavefront adapter", "startup message")
	flags.AddGoFlagSet(flag.CommandLine) // make sure we get the glog flags
	flags.Parse(os.Args)

	wavefrontProvider := cmd.makeProviderOrDie()
	cmd.WithCustomMetrics(wavefrontProvider)
	cmd.WithExternalMetrics(wavefrontProvider)

	glog.Infof("%s version: %s commit tip: %s", cmd.Message, version, commit)
	if err := cmd.Run(wait.NeverStop); err != nil {
		glog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}
