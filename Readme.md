# Metrics Exporter

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

Metrics Exporter is a Go library that provides functionality to interact with Prometheus endpoints and exposes metrics in the text-based exposition format. It simplifies the process of collecting and selectively exposing metrics for monitoring purposes. For example, if there is an usecase of exposing subset of metrics, prometheus metrics exporter enables us to filter the metrics using prometheus query endpoint, and exposes the metrics in text based exposition format.

## Features

- Seamless integration with Prometheus endpoints
- Automatic exposition of metrics in the text-based format

## Installation

To install Your Library Name, use the following command:

```shell
go get -u github.com/milind12/metrics_exporter
```

## Getting Started

To get started with this project, follow these steps:

1 . Import the package into your Go code:

```
import "github.com/milind12/metrics_exporter"
```

2.  Following code example demonstrates selectively exposing metrics which have label `job` with value `demo-app-metrics`

```
    promInstance, err := NewPrometheusInstance("localhost", "9090", map[string]string{"job": "demo-app-metrics"})
	if err != nil {
		t.Errorf("failed")
	}
	exportedMetrics:= promInstance.ExportMetrics()
	fmt.Println(exportedMetrics)
```
