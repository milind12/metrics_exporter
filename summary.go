package metrics_exporter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Quantile struct {
	Qvalue string
	Value  string
}

type SummaryMetric struct {
	Label     string
	Quantiles []Quantile
	Sum       string
	Count     string
}

type SummaryMetrics struct {
	Name             string
	Help             string
	Type             string
	SummaryMetricMap []SummaryMetric
}

func NewSummaryMetrics(Name string, Help string, Type string) *SummaryMetrics {
	return &SummaryMetrics{Name: Name, Help: Help, Type: Type}
}

func (sm *SummaryMetrics) ParseMetricMap(MetricMap map[string]string, Value string) {
	labels := make([]string, 0, len(MetricMap))
	metricName := MetricMap["__name__"]
	for k, v := range MetricMap {
		if k == "__name__" || k == "job" || k == "instance" || k == "quantile" {
			continue
		}
		labels = append(labels, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(labels)
	sortedLabels := strings.Join(labels, ",")
	summaryMetric := sm.FindSummaryMetricByLabel(sortedLabels)

	if strings.HasSuffix(metricName, "_sum") {
		summaryMetric.Sum = Value
	} else if strings.HasSuffix(metricName, "_count") {
		summaryMetric.Count = Value
	} else {
		if summaryMetric.Quantiles == nil {
			summaryMetric.Quantiles = []Quantile{}
		}
		quantile := Quantile{Qvalue: MetricMap["quantile"], Value: Value}
		summaryMetric.Quantiles = append(summaryMetric.Quantiles, quantile)
	}
}

func (sm *SummaryMetrics) Sort() {
	sort.Slice(sm.SummaryMetricMap, func(i, j int) bool {
		return sm.SummaryMetricMap[i].Label < sm.SummaryMetricMap[j].Label
	})

	for _, summaryMetric := range sm.SummaryMetricMap {
		sort.Slice(summaryMetric.Quantiles, func(i, j int) bool {
			ithLen, _ := strconv.ParseFloat(summaryMetric.Quantiles[i].Qvalue, 64)
			jthLen, _ := strconv.ParseFloat(summaryMetric.Quantiles[j].Qvalue, 64)
			return ithLen < jthLen
		})
	}
}

func (sm *SummaryMetrics) InitSummaryMetricByLabel(label string) *SummaryMetric {
	bucketMetrics := &SummaryMetric{Label: label}
	if sm.SummaryMetricMap == nil {
		sm.SummaryMetricMap = []SummaryMetric{}
	}
	sm.SummaryMetricMap = append(sm.SummaryMetricMap, *bucketMetrics)
	return &sm.SummaryMetricMap[len(sm.SummaryMetricMap)-1]
}

func (sm *SummaryMetrics) FindSummaryMetricByLabel(label string) *SummaryMetric {

	for i, _ := range sm.SummaryMetricMap {
		if sm.SummaryMetricMap[i].Label == label {
			return &sm.SummaryMetricMap[i]
		}
	}
	return sm.InitSummaryMetricByLabel(label)
}

func (sm *SummaryMetrics) Print() string {
	sm.Sort()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# HELP %s %s\n", sm.Name, sm.Help))
	builder.WriteString(fmt.Sprintf("# TYPE %s %s\n", sm.Name, sm.Type))
	printLabels := false
	if len(sm.SummaryMetricMap) > 1 {
		printLabels = true
	}
	for _, summaryMetric := range sm.SummaryMetricMap {
		for _, Quantile := range summaryMetric.Quantiles {
			if summaryMetric.Label != "" {
				builder.WriteString(fmt.Sprintf("%s{%s,quantile=\"%s\"} %s\n", sm.Name, summaryMetric.Label, Quantile.Qvalue, Quantile.Value))
			} else {
				builder.WriteString(fmt.Sprintf("%s{quantile=\"%s\"} %s\n", sm.Name, Quantile.Qvalue, Quantile.Value))
			}
		}
		if printLabels {
			builder.WriteString(fmt.Sprintf("%s{%s} %s\n", sm.Name+"_sum", summaryMetric.Label, summaryMetric.Sum))
			builder.WriteString(fmt.Sprintf("%s{%s} %s\n", sm.Name+"_count", summaryMetric.Label, summaryMetric.Count))

		} else {
			builder.WriteString(fmt.Sprintf("%s %s\n", sm.Name+"_sum", summaryMetric.Sum))
			builder.WriteString(fmt.Sprintf("%s %s\n", sm.Name+"_count", summaryMetric.Count))
		}
	}
	return builder.String()
}
