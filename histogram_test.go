package metrics_exporter

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"
)

var hjsonStr string = `{
	"status": "success",
	"data": {
			"resultType": "vector",
			"result": [
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "+Inf",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "+Inf",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "0.01",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "0.01",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "0.1",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "0.1",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "0.5",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "0.5",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "1",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "1",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "10",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "10",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "2.5",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "2.5",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "5",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "5",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "50",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_bucket",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"le": "50",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_count",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_count",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"2"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_sum",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"request_path": "/hello",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"0.000012508"
							]
					},
					{
							"metric": {
									"__name__": "example_app_http_request_duration_seconds_sum",
									"host_name": "localhost:3333",
									"http_method": "GET",
									"instance": "localhost:8000",
									"job": "demo-app-metrics",
									"request_path": "/hello/{id}",
									"status_code": "200"
							},
							"value": [
									1684491411.535,
									"0.000009297"
							]
					}
			 ]
		}
}`

func TestHistogram(t *testing.T) {
	var prometheusQueryAPIResponse PrometheusQueryAPIResponse
	json.Unmarshal([]byte(hjsonStr), &prometheusQueryAPIResponse)
	hm := NewHistogramMetrics("example_app_http_request_duration_seconds", "example seconds", "histogram")
	for _, result := range prometheusQueryAPIResponse.Data.Result {
		hm.ParseMetricMap(result.Metric, result.Value[1].(string))
	}

	reader := strings.NewReader(hm.Print())
	fmt.Println(hm.Print())
	parser := expfmt.TextParser{}
	_, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		t.Errorf("Error happened" + err.Error())
	}

}
