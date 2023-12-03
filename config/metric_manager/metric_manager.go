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
package metric_manager

import (
	"errors"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
)

var (
	instance *MetricManager
	once     sync.Once
)

func NewMetricManager() *MetricManager {
	once.Do(func() {
		instance = &MetricManager{
			metrics: make(map[string]prometheus.Collector),
		}
	})
	return instance
}
func (m *MetricManager) CreateMetric(metricType MetricType, name, help string, labelKeys LabelsKeys) {
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

func (m *MetricManager) IncrementCounter(name string, labelValues Labels) error {
	if metric, ok := m.metrics[m.makeFullNameMetric(name)]; ok {
		if counterVec, ok := metric.(*prometheus.CounterVec); ok {
			counterVec.With(prometheus.Labels(labelValues)).Inc() // Convertendo explicitamente aqui
			return nil
		}
	}
	return errors.New("metric not found or is not a Counter")
}
func (m *MetricManager) SetGauge(name string, value float64, labelValues Labels) error {
	if metric, ok := m.metrics[m.makeFullNameMetric(name)]; ok {
		if gaugeVec, ok := metric.(*prometheus.GaugeVec); ok {
			gaugeVec.With(prometheus.Labels(labelValues)).Set(value)
			return nil
		}
	}
	return errors.New("metric not found or is not a Gauge")
}

func (m *MetricManager) ObserveHistogram(name string, value float64, labelValues Labels) error {
	if metric, ok := m.metrics[m.makeFullNameMetric(name)]; ok {
		if histogramVec, ok := metric.(*prometheus.HistogramVec); ok {
			histogramVec.With(prometheus.Labels(labelValues)).Observe(value)
			return nil
		}
	}
	return errors.New("metric not found or is not a Histogram")
}

func (m *MetricManager) ObserveSummary(name string, value float64, labelValues Labels) error {
	if metric, ok := m.metrics[name]; ok {
		if summaryVec, ok := metric.(*prometheus.SummaryVec); ok {
			summaryVec.With(prometheus.Labels(labelValues)).Observe(value)
			return nil
		}
	}
	return errors.New("metric not found or is not a Summary")
}

func (m *MetricManager) makeFullNameMetric(name string) string {
	return env.GetEnv("PROMETHEUS_PREFIX", "") + "_" + name
}

func (m *MetricManager) safeRegister(fullName string, collector prometheus.Collector) {
	err := prometheus.Register(collector)
	if err != nil {
		if _, dup := err.(prometheus.AlreadyRegisteredError); !dup {
			log.Printf("Failed to register metric: %v", err)
		}
	} else {
		m.metrics[fullName] = collector
	}
}
