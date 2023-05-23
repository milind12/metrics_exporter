package metrics_exporter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Le struct {
	Length string
	Value  string
}

type HistogramMetric struct {
	Label    string
	BucketLe []Le
	Sum      string
	Count    string
}

type HistogramMetrics struct {
	Name             string
	Help             string
	Type             string
	BucketMetricsMap []HistogramMetric
}

func NewHistogramMetrics(Name string, Help string, Type string) *HistogramMetrics {
	return &HistogramMetrics{Name: Name, Help: Help, Type: Type}
}

func (hm *HistogramMetrics) ParseMetricMap(MetricMap map[string]string, Value string) {
	labels := make([]string, 0, len(MetricMap))
	metricName := MetricMap["__name__"]
	for k, v := range MetricMap {
		if k == "__name__" || k == "job" || k == "instance" || k == "le" {
			continue
		}
		labels = append(labels, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(labels)
	sortedLabels := strings.Join(labels, ",")
	bucketMetrics := hm.FindBucketMetricsByLabel(sortedLabels)

	if strings.HasSuffix(metricName, "_sum") {
		bucketMetrics.Sum = Value
	} else if strings.HasSuffix(metricName, "_count") {
		bucketMetrics.Count = Value
	} else if strings.HasSuffix(metricName, "_bucket") {
		if bucketMetrics.BucketLe == nil {
			bucketMetrics.BucketLe = []Le{}
		}
		le := Le{Length: MetricMap["le"], Value: Value}
		bucketMetrics.BucketLe = append(bucketMetrics.BucketLe, le)
	}
}

func (hm *HistogramMetrics) Sort() {
	sort.Slice(hm.BucketMetricsMap, func(i, j int) bool {
		return hm.BucketMetricsMap[i].Label < hm.BucketMetricsMap[j].Label
	})

	for _, bm := range hm.BucketMetricsMap {
		sort.Slice(bm.BucketLe, func(i, j int) bool {
			if bm.BucketLe[i].Length == "-Inf" {
				return true
			} else if bm.BucketLe[j].Length == "-Inf" {
				return false
			} else if bm.BucketLe[i].Length == "+Inf" {
				return false
			} else if bm.BucketLe[j].Length == "+Inf" {
				return true
			}
			ithLen, _ := strconv.ParseFloat(bm.BucketLe[i].Length, 64)
			jthLen, _ := strconv.ParseFloat(bm.BucketLe[j].Length, 64)
			return ithLen < jthLen
		})
	}
}

func (hm *HistogramMetrics) InitBucketMetricsByLabel(label string) *HistogramMetric {
	bucketMetrics := &HistogramMetric{Label: label}
	if hm.BucketMetricsMap == nil {
		hm.BucketMetricsMap = []HistogramMetric{}
	}
	hm.BucketMetricsMap = append(hm.BucketMetricsMap, *bucketMetrics)
	return &hm.BucketMetricsMap[len(hm.BucketMetricsMap)-1]
}

func (hm *HistogramMetrics) FindBucketMetricsByLabel(label string) *HistogramMetric {

	for i, _ := range hm.BucketMetricsMap {
		if hm.BucketMetricsMap[i].Label == label {
			return &hm.BucketMetricsMap[i]
		}
	}
	return hm.InitBucketMetricsByLabel(label)
}

func (hm *HistogramMetrics) Print() string {
	hm.Sort()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# HELP %s %s\n", hm.Name, hm.Help))
	builder.WriteString(fmt.Sprintf("# TYPE %s %s\n", hm.Name, hm.Type))
	printLabels := false
	if len(hm.BucketMetricsMap) > 1 {
		printLabels = true
	}
	for _, bucket := range hm.BucketMetricsMap {
		for _, le := range bucket.BucketLe {
			if bucket.Label != "" {
				builder.WriteString(fmt.Sprintf("%s{%s,le=\"%s\"} %s\n", hm.Name+"_bucket", bucket.Label, le.Length, le.Value))
			} else {
				builder.WriteString(fmt.Sprintf("%s{le=\"%s\"} %s\n", hm.Name+"_bucket", le.Length, le.Value))
			}
		}
		if printLabels {
			builder.WriteString(fmt.Sprintf("%s{%s} %s\n", hm.Name+"_sum", bucket.Label, bucket.Sum))
			builder.WriteString(fmt.Sprintf("%s{%s} %s\n", hm.Name+"_count", bucket.Label, bucket.Count))

		} else {
			builder.WriteString(fmt.Sprintf("%s %s\n", hm.Name+"_sum", bucket.Sum))
			builder.WriteString(fmt.Sprintf("%s %s\n", hm.Name+"_count", bucket.Count))
		}
	}
	return builder.String()
}
