package metrics_exporter

import (
	"strings"
)

type PrometheusQueryAPIResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type Metadata []struct {
	Type string `json:"type"`
	Help string `json:"help"`
	Unit string `json:"unit"`
}

type MetricsOrder struct {
	mo []string
}

type PrometheusMetadataAPIResponse struct {
	Status string `json:"status"`
	Data   map[string][]struct {
		Type string `json:"type"`
		Help string `json:"help"`
		Unit string `json:"unit"`
	} `json:"data"`
}

type IMetric interface {
	Print() string
	ParseMetricMap(MetricMap map[string]string, Value string)
}

type MetricsMap map[string]IMetric

func getPossibleMetricName(mname string) string {
	lastUnderscoreIndex := strings.LastIndex(mname, "_")

	if lastUnderscoreIndex != -1 {
		return mname[:lastUnderscoreIndex]
	}
	return ""
}

func initMetric(Name string, Help string, Type string) IMetric {
	if Type == "histogram" {
		return NewHistogramMetrics(Name, Help, Type)
	} else if Type == "counter" {
		return NewCountersMetrics(Name, Help, Type)
	} else if Type == "gauge" {
		return NewGaugesMetrics(Name, Help, Type)
	}
	return NewSummaryMetrics(Name, Help, Type)
}

func getMetric(mm MetricsMap, mn string, mtd Metadata, metricsorder *MetricsOrder) IMetric {
	metric, ok := mm[mn]
	if !ok {
		mm[mn] = initMetric(mn, mtd[0].Help, mtd[0].Type)
		metricsorder.mo = append(metricsorder.mo, mn)
		return mm[mn]
	}
	return metric
}

func ExportMetrics(promQueryResponse PrometheusQueryAPIResponse, promMetadataResponse PrometheusMetadataAPIResponse) string {
	mm := MetricsMap{}
	metricsorder := MetricsOrder{mo: []string{}}
	var builder strings.Builder
	for _, result := range promQueryResponse.Data.Result {
		mn := result.Metric["__name__"]
		metadata, ok := promMetadataResponse.Data[mn]
		if !ok {
			mn = getPossibleMetricName(mn)
			metadata, ok = promMetadataResponse.Data[mn]
		}
		if ok {
			metric := getMetric(mm, mn, metadata, &metricsorder)
			metric.ParseMetricMap(result.Metric, result.Value[1].(string))
		}
	}

	for _, mn := range metricsorder.mo {
		builder.WriteString(mm[mn].Print())
		builder.WriteString("\n")
	}

	return builder.String()
}
