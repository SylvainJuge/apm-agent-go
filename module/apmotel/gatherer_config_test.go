// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package apmotel // import "go.elastic.co/apm/module/apmotel/v2"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestNewGathererConfig(t *testing.T) {
	aggregationSelector := func(metric.InstrumentKind) aggregation.Aggregation { return nil }
	temporalitySelector := func(metric.InstrumentKind) metricdata.Temporality { return metricdata.CumulativeTemporality }

	testCases := []struct {
		name       string
		options    []GathererOption
		wantConfig gathererConfig
	}{
		{
			name:       "Default",
			options:    nil,
			wantConfig: gathererConfig{},
		},
		{
			name: "WithAggregationSelector",
			options: []GathererOption{
				WithAggregationSelector(aggregationSelector),
			},
			wantConfig: gathererConfig{},
		},
		{
			name: "WithTemporalitySelector",
			options: []GathererOption{
				WithTemporalitySelector(temporalitySelector),
			},
			wantConfig: gathererConfig{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newGathererConfig(tt.options...)
			// tested by TestConfigManualReaderOptions
			cfg.aggregation = nil
			cfg.temporality = nil

			assert.Equal(t, tt.wantConfig, cfg)
		})
	}
}

func TestConfigManualReaderOptions(t *testing.T) {
	aggregationSelector := func(metric.InstrumentKind) aggregation.Aggregation { return nil }
	temporalitySelector := func(metric.InstrumentKind) metricdata.Temporality { return metricdata.CumulativeTemporality }

	testCases := []struct {
		name            string
		config          gathererConfig
		wantOptionCount int
	}{
		{
			name:            "Default",
			config:          gathererConfig{},
			wantOptionCount: 0,
		},
		{
			name:            "WithAggregationSelector",
			config:          gathererConfig{aggregation: aggregationSelector},
			wantOptionCount: 1,
		},
		{
			name:            "WithTemporalitySelector",
			config:          gathererConfig{temporality: temporalitySelector},
			wantOptionCount: 1,
		},
		{
			name:            "WithAggregationSelectorWithTemporalitySelector",
			config:          gathererConfig{aggregation: aggregationSelector, temporality: temporalitySelector},
			wantOptionCount: 2,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.config.manualReaderOptions()
			assert.Len(t, opts, tt.wantOptionCount)
		})
	}
}
