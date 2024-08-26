/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metric

import (
	"errors"
	"sync"

	"github.com/caiomarcatti12/nanogo/v1/pkg/env"
	"github.com/caiomarcatti12/nanogo/v1/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	instance *Prometheus
	once     sync.Once
)

type Prometheus struct {
	logger    log.ILog
	namespace string
	metrics   map[string]prometheus.Collector
}

func NewInstancePrometheus(env env.IEnv, logger log.ILog) IMetric {
	logger.Info("Creating instance Prometheus...")
	once.Do(func() {
		instance = &Prometheus{
			logger:    logger,
			namespace: env.GetEnv("PROMETHEUS_PREFIX", "nanogo"),
			metrics:   make(map[string]prometheus.Collector),
		}
	})

	return instance
}

func (m *Prometheus) CreateMetric(metricType MetricType, name, help string, labelKeys LabelsKeys) {
	m.logger.Infof("Creating metric Prometheus %s, %s...", metricType, name)

	fullName := m.makeFullNameMetric(name)

	switch metricType {
	case Counter:
		counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fullName,
			Help: help,
		}, labelKeys)
		m.safeRegister(fullName, counterVec)
	case Gauge:
		gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: fullName,
			Help: help,
		}, labelKeys)
		m.safeRegister(fullName, gaugeVec)
	case Histogram:
		histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: fullName,
			Help: help,
		}, labelKeys)
		m.safeRegister(fullName, histogramVec)
	case Summary:
		summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: fullName,
			Help: help,
		}, labelKeys)
		m.safeRegister(fullName, summaryVec)
	default:
		// Handle invalid metric type
	}
}

func (m *Prometheus) IncrementCounter(name string, labelValues Labels) error {
	if metric, ok := m.metrics[m.makeFullNameMetric(name)]; ok {
		if counterVec, ok := metric.(*prometheus.CounterVec); ok {
			counterVec.With(prometheus.Labels(labelValues)).Inc() // Convertendo explicitamente aqui
			return nil
		}
	}

	m.logger.Errorf("metric %s not found or is not a Counter", name)
	return errors.New("metric not found or is not a Counter")
}

func (m *Prometheus) SetGauge(name string, value float64, labelValues Labels) error {
	if metric, ok := m.metrics[m.makeFullNameMetric(name)]; ok {
		if gaugeVec, ok := metric.(*prometheus.GaugeVec); ok {
			gaugeVec.With(prometheus.Labels(labelValues)).Set(value)
			return nil
		}
	}

	m.logger.Errorf("metric %s not found or is not a Gauge", name)
	return errors.New("metric not found or is not a Gauge")
}

func (m *Prometheus) ObserveHistogram(name string, value float64, labelValues Labels) error {
	if metric, ok := m.metrics[m.makeFullNameMetric(name)]; ok {
		if histogramVec, ok := metric.(*prometheus.HistogramVec); ok {
			histogramVec.With(prometheus.Labels(labelValues)).Observe(value)
			return nil
		}
	}

	m.logger.Errorf("metric %s not found or is not a Histogram", name)
	return errors.New("metric not found or is not a Histogram")
}

func (m *Prometheus) ObserveSummary(name string, value float64, labelValues Labels) error {
	if metric, ok := m.metrics[name]; ok {
		if summaryVec, ok := metric.(*prometheus.SummaryVec); ok {
			summaryVec.With(prometheus.Labels(labelValues)).Observe(value)
			return nil
		}
	}

	m.logger.Errorf("metric %s not found or is not a Summary", name)
	return errors.New("metric not found or is not a Summary")
}

func (m *Prometheus) makeFullNameMetric(name string) string {
	return m.namespace + "_" + name
}

func (m *Prometheus) safeRegister(fullName string, collector prometheus.Collector) {
	err := prometheus.Register(collector)
	if err != nil {
		if _, dup := err.(prometheus.AlreadyRegisteredError); !dup {
			m.logger.Error("Failed to register collector %s: %s", fullName, err)
		}
	} else {
		m.metrics[fullName] = collector
	}
}
