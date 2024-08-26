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
	"fmt"

	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

type MetricType string

// Metricser define a interface for our metrics manager.
type IMetric interface {
	CreateMetric(metricType MetricType, name, help string, labelKeys LabelsKeys)
	IncrementCounter(name string, labelValues Labels) error
	SetGauge(name string, value float64, labelValues Labels) error
	ObserveHistogram(name string, value float64, labelValues Labels) error
	ObserveSummary(name string, value float64, labelValues Labels) error
}

type Labels map[string]string
type LabelsKeys []string

const (
	Counter   MetricType = "Counter"
	Gauge     MetricType = "Gauge"
	Histogram MetricType = "Histogram"
	Summary   MetricType = "Summary"
)

func Factory(env env.IEnv, logger log.ILog) IMetric {
	logger.Info("Creating metric provider...")

	provider := env.GetEnv("METRIC_PROVIDER", "PROMETHEUS")

	switch provider {
	case "PROMETHEUS":
		instance := NewInstancePrometheus(env, logger)

		return instance
	default:
		panic(fmt.Errorf("metric provider %s not found", provider))
	}
}
