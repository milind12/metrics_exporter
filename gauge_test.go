package metrics_exporter

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"
)

var gjsonStr string = `{
	"status": "success",
	"data": {
	  "resultType": "vector",
	  "result": [
		{
			"metric": {
				"__name__": "go_goroutines",
				"instance": "localhost:8000",
				"job": "demo-app-metrics"
			},
			"value": [
				1684735712.333,
				"11"
			]
		}
	  ]
	}
  }`

func TestGauge(t *testing.T) {
	var prometheusQueryAPIResponse PrometheusQueryAPIResponse
	json.Unmarshal([]byte(gjsonStr), &prometheusQueryAPIResponse)
	gm := NewCountersMetrics("go_goroutines", "example Help", "gauge")
	for _, result := range prometheusQueryAPIResponse.Data.Result {
		gm.ParseMetricMap(result.Metric, result.Value[1].(string))
	}
	reader := strings.NewReader(gm.Print())
	fmt.Println(gm.Print())
	parser := expfmt.TextParser{}
	_, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		t.Errorf("Error happened" + err.Error())
	}

}
