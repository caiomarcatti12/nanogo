package metric

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakeEnv is a simple implementation of env.IEnv for tests.
type fakeEnv struct {
	values map[string]string
}

func (f *fakeEnv) GetEnv(key string, defaultValue ...string) string {
	if v, ok := f.values[key]; ok && v != "" {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (f *fakeEnv) GetEnvBool(key string, defaultValue ...string) bool {
	val := f.GetEnv(key, defaultValue...)
	b, _ := strconv.ParseBool(val)
	return b
}

// fakeLogger is a no-op logger used in tests.
type fakeLogger struct{}

func (fakeLogger) Fatal(message string, args ...interface{})   {}
func (fakeLogger) Debug(message string, args ...interface{})   {}
func (fakeLogger) Info(message string, args ...interface{})    {}
func (fakeLogger) Error(message string, args ...interface{})   {}
func (fakeLogger) Warning(message string, args ...interface{}) {}
func (fakeLogger) Trace(message string, args ...interface{})   {}

func (fakeLogger) Fatalf(message string, args ...interface{})   {}
func (fakeLogger) Debugf(message string, args ...interface{})   {}
func (fakeLogger) Infof(message string, args ...interface{})    {}
func (fakeLogger) Errorf(message string, args ...interface{})   {}
func (fakeLogger) Warningf(message string, args ...interface{}) {}
func (fakeLogger) Tracef(message string, args ...interface{})   {}

func TestCreateMetric_UsesNamespacePrefix(t *testing.T) {
	env := &fakeEnv{values: map[string]string{"PROMETHEUS_PREFIX": "testns"}}
	logger := fakeLogger{}
	m := NewInstancePrometheus(env, logger).(*Prometheus)

	m.CreateMetric(Summary, "requests_total", "", LabelsKeys{})

	_, ok := m.metrics["testns_requests_total"]
	assert.True(t, ok, "metric should be stored with namespace prefix")
}

func TestObserveSummary_UsesFullName(t *testing.T) {
	env := &fakeEnv{values: map[string]string{"PROMETHEUS_PREFIX": "testns"}}
	logger := fakeLogger{}
	m := NewInstancePrometheus(env, logger).(*Prometheus)

	m.CreateMetric(Summary, "duration_seconds", "", LabelsKeys{})
	err := m.ObserveSummary("duration_seconds", 1.0, Labels{})

	assert.NoError(t, err)
}
