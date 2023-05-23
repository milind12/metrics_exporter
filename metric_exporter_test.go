package metrics_exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/prometheus/common/expfmt"
)

func fetchFromFile(filePath string) []byte {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read JSON file: %v\n", err)
		return []byte{}
	}
	return jsonBytes
}

func TestMetricsExporter(t *testing.T) {
	queryResponse := fetchFromFile("query.json")
	metadataResponse := fetchFromFile("metadata.json")
	var prometheusQueryAPIResponse PrometheusQueryAPIResponse
	json.Unmarshal(queryResponse, &prometheusQueryAPIResponse)
	var prometheusMetadataAPIResponse PrometheusMetadataAPIResponse
	json.Unmarshal(metadataResponse, &prometheusMetadataAPIResponse)
	exportedMetrics := ExportMetrics(prometheusQueryAPIResponse, prometheusMetadataAPIResponse)
	fmt.Println(exportedMetrics)
	reader := strings.NewReader(exportedMetrics)
	parser := expfmt.TextParser{}
	_, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		t.Errorf("Error happened" + err.Error())
	}
}
