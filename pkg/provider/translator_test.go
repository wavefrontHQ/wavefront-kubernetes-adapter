package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitMetric(t *testing.T) {
	res, metric := splitMetric("kubernetes", "kubernetes.pod.cpu.usage")
	assert.Equal(t, res, "pod")
	assert.Equal(t, metric, "cpu.usage")

	res, metric = splitMetric("pks.kubernetes", "pks.kubernetes.pod.cpu.usage")
	assert.Equal(t, res, "pod")
	assert.Equal(t, metric, "cpu.usage")
}
