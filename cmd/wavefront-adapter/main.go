// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/component-base/logs"

	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	customprovider "sigs.k8s.io/custom-metrics-apiserver/pkg/provider"

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
	// Wavefront client timeout
	APIClientTimeout time.Duration
	// Wavefront Server URL of the form https://INSTANCE.wavefront.com
	WavefrontServerURL string
	// Wavefront API token with permissions to query points
	WavefrontAPIToken string
	// The prefix for custom kubernetes metrics in Wavefront
	CustomMetricPrefix string
	// The file containing the metrics discovery configuration
	AdapterConfigFile string
	// The log level
	LogLevel string
}

func (a *WavefrontAdapter) makeProviderOrDie() customprovider.MetricsProvider {
	conf, err := a.ClientConfig()
	if err != nil {
		log.Fatalf("error getting kube config: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatalf("error creating kube client: %v", err)
	}

	dynClient, err := a.DynamicClient()
	if err != nil {
		log.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		log.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	waveURL, err := url.Parse(a.WavefrontServerURL)
	if err != nil {
		log.Fatalf("unable to parse wavefront url: %v", err)
	}
	waveClient := client.NewWavefrontClient(waveURL, a.WavefrontAPIToken, a.APIClientTimeout)

	metricsProvider, runnable := provider.NewWavefrontProvider(provider.WavefrontProviderConfig{
		DynClient:    dynClient,
		KubeClient:   kubeClient,
		Mapper:       mapper,
		WaveClient:   waveClient,
		Prefix:       strings.Trim(a.CustomMetricPrefix, "."),
		ListInterval: a.MetricsRelistInterval,
		ExternalCfg:  a.AdapterConfigFile,
	})
	runnable.RunUntil(wait.NeverStop)
	return metricsProvider
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)

	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	cmd := &WavefrontAdapter{
		CustomMetricPrefix:    "kubernetes",
		MetricsRelistInterval: 10 * time.Minute,
		APIClientTimeout:      10 * time.Second,
	}
	cmd.Name = "wavefront-custom-metrics-adapter"
	flags := cmd.Flags()
	flags.DurationVar(&cmd.MetricsRelistInterval, "metrics-relist-interval", cmd.MetricsRelistInterval, ""+
		"Interval at which to fetch the list of custom metric names from Operations for Applications.")
	flags.DurationVar(&cmd.APIClientTimeout, "api-client-timeout", cmd.APIClientTimeout, ""+
		"Client timeout to Operations for Applications.")
	flags.StringVar(&cmd.WavefrontServerURL, "wavefront-url", "",
		"Wavefront URL in the format https://YOUR_INSTANCE.wavefront.com")
	flags.StringVar(&cmd.WavefrontAPIToken, "wavefront-token", "",
		"Wavefront API token with permissions to query for points.")
	flags.StringVar(&cmd.CustomMetricPrefix, "wavefront-metric-prefix", cmd.CustomMetricPrefix,
		"Metrics under this prefix are exposed in the custom metrics API.")
	flags.StringVar(&cmd.AdapterConfigFile, "external-metrics-config", "",
		"Configuration file for driving external metrics API.")
	flags.StringVar(&cmd.LogLevel, "log-level", "info", "One of info, debug or trace.")
	flags.StringVar(&cmd.Message, "msg", "starting wavefront adapter", "startup message")
	flags.AddGoFlagSet(flag.CommandLine) // make sure we get the glog flags
	flags.Parse(os.Args)

	switch cmd.LogLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	}

	wavefrontProvider := cmd.makeProviderOrDie()
	cmd.WithCustomMetrics(wavefrontProvider)
	cmd.WithExternalMetrics(wavefrontProvider)

	log.Infof("%s version: %s commit tip: %s", cmd.Message, version, commit)
	if err := cmd.Run(wait.NeverStop); err != nil {
		log.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}
