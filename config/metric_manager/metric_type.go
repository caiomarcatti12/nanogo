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
	"github.com/prometheus/client_golang/prometheus"
)

type MetricType int

// Metricser define a interface for our metrics manager.
type MetricManagerInterface interface {
	CreateMetric(metricType MetricType, name, help string, labelKeys LabelsKeys)
	IncrementCounter(name string, labelValues Labels) error
	SetGauge(name string, value float64, labelValues Labels) error
	ObserveHistogram(name string, value float64, labelValues Labels) error
	ObserveSummary(name string, value float64, labelValues Labels) error
}

type MetricManager struct {
	metrics map[string]prometheus.Collector
}

type Labels map[string]string
type LabelsKeys []string

const (
	Counter MetricType = iota
	Gauge
	Histogram
	Summary
)
