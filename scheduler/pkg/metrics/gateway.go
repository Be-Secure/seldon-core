/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed by
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/

package metrics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Keep next line as used in docs
// Docs Start Metrics
//
// The aggregate metrics exist for efficiency, as the summation can be
// very slow in Prometheus when many pipelines exist.
const (
	// Histograms do no include model label for efficiency
	pipelineHistogramName = "seldon_pipeline_infer_api_seconds"
	// We use base infer counters to store core metrics per pipeline
	pipelineInferCounterName                 = "seldon_pipeline_infer_total"
	pipelineInferLatencyCounterName          = "seldon_pipeline_infer_seconds_total"
	pipelineAggregateInferCounterName        = "seldon_pipeline_aggregate_infer_total"
	pipelineAggregateInferLatencyCounterName = "seldon_pipeline_aggregate_infer_seconds_total"
)

// Docs End Metrics
// Keep above line as used in docs
// Metric labels
const (
	SeldonPipelineMetric = "pipeline"
)

// TODO Revisit splitting this interface as metric handling matures
type PipelineMetricsHandler interface {
	AddPipelineInferMetrics(pipelineName string, method string, elapsedTime float64, code string)
}

type PrometheusPipelineMetrics struct {
	serverName string
	logger     log.FieldLogger
	// Model metrics
	pipelineHistogram                    *prometheus.HistogramVec
	pipelineInferCounter                 *prometheus.CounterVec
	pipelineInferLatencyCounter          *prometheus.CounterVec
	pipelineAggregateInferCounter        *prometheus.CounterVec
	pipelineAggregateInferLatencyCounter *prometheus.CounterVec

	server *http.Server
}

func NewPrometheusPipelineMetrics(logger log.FieldLogger) (*PrometheusPipelineMetrics, error) {
	histogram, err := createPipelineHistogram()
	if err != nil {
		return nil, err
	}

	inferCounter, err := createPipelineInferCounter()
	if err != nil {
		return nil, err
	}

	inferLatencyCounter, err := createPipelineInferLatencyCounter()
	if err != nil {
		return nil, err
	}

	aggregateInferCounter, err := createPipelineAggregateInferCounter()
	if err != nil {

		return nil, err
	}

	aggregateInferLatencyCounter, err := createPipelineAggregateInferLatencyCounter()
	if err != nil {
		return nil, err
	}

	return &PrometheusPipelineMetrics{
		serverName:                           "pipeline-gateway",
		logger:                               logger.WithField("source", "PrometheusMetrics"),
		pipelineHistogram:                    histogram,
		pipelineInferCounter:                 inferCounter,
		pipelineInferLatencyCounter:          inferLatencyCounter,
		pipelineAggregateInferCounter:        aggregateInferCounter,
		pipelineAggregateInferLatencyCounter: aggregateInferLatencyCounter,
	}, nil
}

func createPipelineHistogram() (*prometheus.HistogramVec, error) {
	//TODO add method for rest/grpc
	labelNames := []string{SeldonServerMetric, SeldonPipelineMetric, MethodMetric, CodeMetric}

	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    pipelineHistogramName,
			Help:    "A histogram of latencies for pipeline gateway",
			Buckets: DefaultHistogramBuckets,
		},
		labelNames,
	)
	err := prometheus.Register(histogram)
	if err != nil {
		if e, ok := err.(prometheus.AlreadyRegisteredError); ok {
			histogram = e.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			return nil, err
		}
	}
	return histogram, nil
}

func createPipelineInferCounter() (*prometheus.CounterVec, error) {
	labelNames := []string{SeldonServerMetric, SeldonPipelineMetric, MethodTypeMetric, CodeMetric}
	return createCounterVec(
		pipelineInferCounterName,
		"A count of pipeline gateway calls",
		labelNames,
	)
}

func createPipelineInferLatencyCounter() (*prometheus.CounterVec, error) {
	labelNames := []string{SeldonServerMetric, SeldonPipelineMetric, MethodTypeMetric, CodeMetric}
	return createCounterVec(
		pipelineInferLatencyCounterName,
		"A sum of pipeline gateway call latencies",
		labelNames,
	)
}

func createPipelineAggregateInferCounter() (*prometheus.CounterVec, error) {
	labelNames := []string{SeldonServerMetric, MethodTypeMetric}
	return createCounterVec(
		pipelineAggregateInferCounterName,
		"A count of pipeline gateway calls (aggregate)",
		labelNames,
	)
}

func createPipelineAggregateInferLatencyCounter() (*prometheus.CounterVec, error) {
	labelNames := []string{SeldonServerMetric, MethodTypeMetric}
	return createCounterVec(
		pipelineAggregateInferLatencyCounterName,
		"A sum of pipeline gateway call latencies (aggregate)",
		labelNames,
	)
}

func (pm *PrometheusPipelineMetrics) HttpCodeToString(code int) string {
	return fmt.Sprintf("%d", code)
}

func (pm *PrometheusPipelineMetrics) AddPipelineInferMetrics(pipelineName string, method string, latency float64, code string) {
	pm.addInferLatency(pipelineName, method, latency, code)
	pm.addInferCount(pipelineName, method, code)
}

func (pm *PrometheusPipelineMetrics) addInferCount(pipelineName, method string, code string) {
	pm.pipelineInferCounter.With(prometheus.Labels{
		SeldonPipelineMetric: pipelineName,
		SeldonServerMetric:   pm.serverName,
		MethodTypeMetric:     method,
		CodeMetric:           code,
	}).Inc()
	pm.pipelineAggregateInferCounter.With(prometheus.Labels{
		SeldonServerMetric: pm.serverName,
		MethodTypeMetric:   method,
	}).Inc()
}

func (pm *PrometheusPipelineMetrics) addInferLatency(pipelineName, method string, latency float64, code string) {
	pm.pipelineInferLatencyCounter.With(prometheus.Labels{
		SeldonPipelineMetric: pipelineName,
		SeldonServerMetric:   pm.serverName,
		MethodTypeMetric:     method,
		CodeMetric:           code,
	}).Add(latency)
	pm.pipelineAggregateInferLatencyCounter.With(prometheus.Labels{
		SeldonServerMetric: pm.serverName,
		MethodTypeMetric:   method,
	}).Add(latency)
	pm.pipelineHistogram.WithLabelValues(pm.serverName, pipelineName, method, code).Observe(latency)
}

func (pm *PrometheusPipelineMetrics) Start(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	pm.server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	pm.logger.Infof("Starting metrics server on port %d", port)
	return pm.server.ListenAndServe()
}

func (pm *PrometheusPipelineMetrics) Stop() error {
	return pm.server.Shutdown(context.Background())
}
