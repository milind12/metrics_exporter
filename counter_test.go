package metrics_exporter

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"
)

var cjsonStr string = `{
	"status": "success",
	"data": {
	  "resultType": "vector",
	  "result": [
		{
			"metric": {
				"__name__": "example_app_http_requests_total",
				"host_name": "localhost:3333",
				"http_method": "GET",
				"instance": "localhost:8000",
				"job": "demo-app-metrics",
				"request_path": "/hello/{id}",
				"status_code": "200"
			},
			"value": [
				1684735712.333,
				"10"
			]
		},{
			"metric": {
				"__name__": "example_app_http_requests_total",
				"host_name": "localhost:3333",
				"http_method": "GET",
				"instance": "localhost:8000",
				"job": "demo-app-metrics",
				"request_path": "/hello",
				"status_code": "200"
			},
			"value": [
				1684735712.333,
				"4"
			]
		}
	  ]
	}
  }`

func TestCounter(t *testing.T) {
	var prometheusQueryAPIResponse PrometheusQueryAPIResponse
	json.Unmarshal([]byte(cjsonStr), &prometheusQueryAPIResponse)
	cm := NewCountersMetrics("example_app_http_requests_total", "example seconds", "counter")
	for _, result := range prometheusQueryAPIResponse.Data.Result {
		cm.ParseMetricMap(result.Metric, result.Value[1].(string))
	}
	reader := strings.NewReader(cm.Print())
	fmt.Println(cm.Print())
	parser := expfmt.TextParser{}
	_, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		t.Errorf("Error happened" + err.Error())
	}

}
