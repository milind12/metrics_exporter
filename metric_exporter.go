package metrics_exporter

import (
	"strings"
)

type Metadata []struct {
	Type string `json:"type"`
	Help string `json:"help"`
	Unit string `json:"unit"`
}

type MetricsOrder struct {
	mo []string
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

func getMetric(metricmap MetricsMap, metricname string, mtd Metadata, metricsorder *MetricsOrder) IMetric {
	metric, ok := metricmap[metricname]
	if !ok {
		metricmap[metricname] = initMetric(metricname, mtd[0].Help, mtd[0].Type)
		metricsorder.mo = append(metricsorder.mo, metricname)
		return metricmap[metricname]
	}
	return metric
}

func exportMetrics(promQueryResponse PrometheusQueryAPIResponse, promMetadataResponse PrometheusMetadataAPIResponse) string {
	metricmap := MetricsMap{}
	metricsorder := MetricsOrder{mo: []string{}}
	var builder strings.Builder
	for _, result := range promQueryResponse.Data.Result {
		metricname := result.Metric["__name__"]
		metadata, ok := promMetadataResponse.Data[metricname]
		if !ok {
			metricname = getPossibleMetricName(metricname)
			metadata, ok = promMetadataResponse.Data[metricname]
		}
		if ok {
			metric := getMetric(metricmap, metricname, metadata, &metricsorder)
			metric.ParseMetricMap(result.Metric, result.Value[1].(string))
		}
	}

	for _, metricname := range metricsorder.mo {
		builder.WriteString(metricmap[metricname].Print())
		builder.WriteString("\n")
	}

	return builder.String()
}
