package metrics_exporter

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"
)

var sjsonStr string = `{
	"status": "success",
	"data": {
	  "resultType": "vector",
	  "result": [
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds",
			"instance": "localhost:8000",
			"job": "demo-app-metrics",
			"quantile": "0"
		  },
		  "value": [
			1684735712.333,
			"0.000023068"
		  ]
		},
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds",
			"instance": "localhost:8000",
			"job": "demo-app-metrics",
			"quantile": "0.25"
		  },
		  "value": [
			1684735712.333,
			"0.000023068"
		  ]
		},
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds",
			"instance": "localhost:8000",
			"job": "demo-app-metrics",
			"quantile": "0.5"
		  },
		  "value": [
			1684735712.333,
			"0.000023068"
		  ]
		},
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds",
			"instance": "localhost:8000",
			"job": "demo-app-metrics",
			"quantile": "0.75"
		  },
		  "value": [
			1684735712.333,
			"0.000023068"
		  ]
		},
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds",
			"instance": "localhost:8000",
			"job": "demo-app-metrics",
			"quantile": "1"
		  },
		  "value": [
			1684735712.333,
			"0.000023068"
		  ]
		},
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds_count",
			"instance": "localhost:8000",
			"job": "demo-app-metrics"
		  },
		  "value": [
			1684735712.333,
			"1"
		  ]
		},
		{
		  "metric": {
			"__name__": "go_gc_duration_seconds_sum",
			"instance": "localhost:8000",
			"job": "demo-app-metrics"
		  },
		  "value": [
			1684735712.333,
			"0.000023068"
		  ]
		}
	  ]
	}
  }`

func TestSummary(t *testing.T) {
	var prometheusQueryAPIResponse PrometheusQueryAPIResponse
	json.Unmarshal([]byte(sjsonStr), &prometheusQueryAPIResponse)
	sm := NewSummaryMetrics("go_gc_duration_seconds", "example seconds", "summary")
	for _, result := range prometheusQueryAPIResponse.Data.Result {
		sm.ParseMetricMap(result.Metric, result.Value[1].(string))
	}
	reader := strings.NewReader(sm.Print())
	fmt.Println(sm.Print())
	parser := expfmt.TextParser{}
	_, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		t.Errorf("Error happened" + err.Error())
	}

}
