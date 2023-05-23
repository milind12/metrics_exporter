package metrics_exporter

import (
	"fmt"
	"sort"
	"strings"
)

type Gauge struct {
	Label string
	Value string
}

type GaugesMetric struct {
	Name   string
	Help   string
	Type   string
	Gauges []Gauge
}

func NewGaugesMetrics(Name string, Help string, Type string) (gm *GaugesMetric) {
	return &GaugesMetric{Name: Name, Help: Help, Type: Type}
}

func (gm *GaugesMetric) ParseMetricMap(MetricMap map[string]string, Value string) {
	labels := make([]string, 0, len(MetricMap))
	for k, v := range MetricMap {
		if k == "__name__" || k == "job" || k == "instance" {
			continue
		}
		labels = append(labels, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(labels)
	sortedLabels := strings.Join(labels, ",")
	gauge := gm.InitCounterByLabel(sortedLabels)
	gauge.Value = Value
}

func (gm *GaugesMetric) InitCounterByLabel(label string) *Gauge {
	counterMetric := Gauge{Label: label}
	if gm.Gauges == nil {
		gm.Gauges = []Gauge{}
	}
	gm.Gauges = append(gm.Gauges, counterMetric)
	return &gm.Gauges[len(gm.Gauges)-1]
}

func (gm *GaugesMetric) Sort() {
	sort.Slice(gm.Gauges, func(i, j int) bool {
		return gm.Gauges[i].Label < gm.Gauges[j].Label
	})
}

func (gm *GaugesMetric) Print() string {
	gm.Sort()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# HELP %s %s\n", gm.Name, gm.Help))
	builder.WriteString(fmt.Sprintf("# TYPE %s %s\n", gm.Name, gm.Type))

	for _, gauge := range gm.Gauges {
		if gauge.Label == "" {
			builder.WriteString(fmt.Sprintf("%s %s\n", gm.Name, gauge.Value))
		} else {
			builder.WriteString(fmt.Sprintf("%s{%s} %s\n", gm.Name, gauge.Label, gauge.Value))
		}
	}
	return builder.String()
}
