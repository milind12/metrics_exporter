package metrics_exporter

import (
	"fmt"
	"sort"
	"strings"
)

type Counter struct {
	Label string
	Value string
}

type CountersMetric struct {
	Name     string
	Help     string
	Type     string
	Counters []Counter
}

func NewCountersMetrics(Name string, Help string, Type string) (cm *CountersMetric) {
	return &CountersMetric{Name: Name, Help: Help, Type: Type}
}

func (cm *CountersMetric) ParseMetricMap(MetricMap map[string]string, Value string) {
	labels := make([]string, 0, len(MetricMap))
	for k, v := range MetricMap {
		if k == "__name__" || k == "job" || k == "instance" {
			continue
		}
		labels = append(labels, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(labels)
	sortedLabels := strings.Join(labels, ",")
	counter := cm.InitCounterByLabel(sortedLabels)
	counter.Value = Value
}

func (cm *CountersMetric) InitCounterByLabel(label string) *Counter {
	counterMetric := Counter{Label: label}
	if cm.Counters == nil {
		cm.Counters = []Counter{}
	}
	cm.Counters = append(cm.Counters, counterMetric)
	return &cm.Counters[len(cm.Counters)-1]
}

func (cm *CountersMetric) Sort() {
	sort.Slice(cm.Counters, func(i, j int) bool {
		return cm.Counters[i].Label < cm.Counters[j].Label
	})
}

func (cm *CountersMetric) Print() string {
	cm.Sort()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# HELP %s %s\n", cm.Name, cm.Help))
	builder.WriteString(fmt.Sprintf("# TYPE %s %s\n", cm.Name, cm.Type))

	for _, counter := range cm.Counters {
		if counter.Label == "" {
			builder.WriteString(fmt.Sprintf("%s %s\n", cm.Name, counter.Value))
		} else {
			builder.WriteString(fmt.Sprintf("%s{%s} %s\n", cm.Name, counter.Label, counter.Value))
		}
	}
	return builder.String()
}
