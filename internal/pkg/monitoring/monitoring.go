// (C) Copyright 2017-2019 Hewlett Packard Enterprise Development LP

package monitoring

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/cloud/internal/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MethodDimension is a dimension for prometheus metrics
	MethodDimension = "method"
	// PathDimension is a dimension for prometheus metrics
	PathDimension = "path"
	// ServiceDimension is a dimension for prometheus metrics
	ServiceDimension = "service"
	// ResultDimension is a dimension for prometheus metrics
	ResultDimension = "result"
	// StatusCodeDimension is a dimension for prometheus metrics
	StatusCodeDimension = "status_code"

	successStatus = "success"
	failureStatus = "failure"
)

var (
	// ExternalReqCountMetric measures the number of external requests from titan
	ExternalReqCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "titan_external_subrequest_count",
			Help: "Number of external requests from titan",
		},
		[]string{MethodDimension, PathDimension, ServiceDimension, ResultDimension, StatusCodeDimension},
	)
	// ExternalReqDurationMetric measures the external requests duration
	ExternalReqDurationMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "titan_external_subrequest_latency",
			Help: "Time taken by external requests made from titans",
		},
		[]string{MethodDimension, PathDimension, ServiceDimension, ResultDimension, StatusCodeDimension},
	)
	// RequestCounter measures the number of incoming requests
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Counts requests by method and path"},
		[]string{MethodDimension, PathDimension, StatusCodeDimension},
	)
	// RequestTimer measures the incoming request latency
	RequestTimer = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total_time",
			Help: "Total time spent on requests by method and path (in seconds)"},
		[]string{MethodDimension, PathDimension, StatusCodeDimension},
	)

	initSync sync.Once
)

// RegisterMetrics register the metrics with prometheus
func RegisterMetrics() {
	initSync.Do(func() {
		prometheus.MustRegister(ExternalReqCountMetric)
		prometheus.MustRegister(ExternalReqDurationMetric)
		prometheus.MustRegister(RequestCounter)
		prometheus.MustRegister(RequestTimer)
	})
}

// CheckRequestReturnStatus returns the status code and response status, whether the request is a success or a failure
// This check is being done in the handlers, but to keep the monitoring separate, checking it here also.
func CheckRequestReturnStatus(okCodes []int, resp *http.Response, err error) (string, string) {
	ResultStatusMap := map[bool]string{
		true:  successStatus,
		false: failureStatus,
	}
	var ok bool

	if resp == nil {
		return "", ResultStatusMap[ok]
	}

	if err != nil {
		return strconv.Itoa(resp.StatusCode), ResultStatusMap[ok]
	}

	// Allow default OkCodes if none explicitly set
	if okCodes == nil {
		okCodes = utils.DefaultOkCodes(resp.Request.Method)
	}
	okCodes = append(okCodes, http.StatusNotFound)
	for _, code := range okCodes {
		if resp.StatusCode == code {
			ok = true
			break
		}
	}
	return strconv.Itoa(resp.StatusCode), ResultStatusMap[ok]
}
