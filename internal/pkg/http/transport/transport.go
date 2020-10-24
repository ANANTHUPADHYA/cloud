//(C) Copyright 2019 Hewlett Packard Enterprise Development LP

package transport

//go:generate mockgen -package transport -destination transport_mock.go net/http RoundTripper

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloud/internal/pkg/monitoring"
)

//LoggerTransport wraps und tripper interface for logging
type LoggerTransport struct {
	tr http.RoundTripper
}

//NewLoggerTransport wraps roundtripper interface around LoggerTransport
func NewLoggerTransport(tr http.RoundTripper) http.RoundTripper {
	return &LoggerTransport{tr: tr}
}

//RoundTrip logs method and path before calling nested round tripper interface
func (trans *LoggerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Infof(req.Context(), "method: [%s], path:[%s]", req.Method, req.URL.String())
	resp, err := trans.tr.RoundTrip(req)
	if resp != nil {
		if resp.StatusCode < 399 {
			log.Infof(req.Context(),
				"method: [%s], path:[%s], Response status: [%d],",
				req.Method, req.URL.String(), resp.StatusCode,
			)
		} else {
			log.Errorf(req.Context(),
				"method: [%s], path:[%s], Response status: [%d],",
				req.Method, req.URL.String(), resp.StatusCode,
			)
		}
	}
	return resp, err
}

//MetricTransport holds round tripper interface
type MetricTransport struct {
	tr http.RoundTripper
}

//NewMetricTransport wraps roundtripper interface around Transport
func NewMetricTransport(tr http.RoundTripper) http.RoundTripper {
	return &MetricTransport{tr: tr}
}

//RoundTrip calls nested round tripper interface
func (trans *MetricTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := trans.tr.RoundTrip(req)
	latency := time.Since(start)

	statusCode, result := monitoring.CheckRequestReturnStatus(nil, resp, err)
	monitoring.RegisterMetrics()
	labels := prometheus.Labels{
		monitoring.MethodDimension:     req.Method,
		monitoring.PathDimension:       req.URL.Path,
		monitoring.ServiceDimension:    req.URL.Host,
		monitoring.ResultDimension:     result,
		monitoring.StatusCodeDimension: statusCode,
	}
	monitoring.ExternalReqCountMetric.With(labels).Inc()
	monitoring.ExternalReqDurationMetric.With(labels).Add(latency.Seconds())

	return resp, err
}
