package metrics_exporter

import (
	"encoding/json"
	"io/ioutil"

	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HTTPResponse struct {
	Response *http.Response
	Error    error
}

type PrometheusExporter struct {
	Host    string
	Port    string
	Filters map[string]string
}

type EmptyFiltersError struct{}

func (e *EmptyFiltersError) Error() string {
	return "Filters Map is empty"
}

func checkFiltersMapNotEmpty(myMap map[string]string) error {
	if len(myMap) == 0 {
		return &EmptyFiltersError{}
	}
	return nil
}

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

type PrometheusMetadataAPIResponse struct {
	Status string `json:"status"`
	Data   map[string][]struct {
		Type string `json:"type"`
		Help string `json:"help"`
		Unit string `json:"unit"`
	} `json:"data"`
}

func makeHTTPRequest(host, port, path string, queryParams map[string]string) *HTTPResponse {
	url := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   path,
	}

	if len(queryParams) > 0 {
		values := url.Query()
		for key, value := range queryParams {
			values.Add(key, value)
		}
		url.RawQuery = values.Encode()
	}

	response, err := http.Get(url.String())
	return &HTTPResponse{Response: response, Error: err}
}

func NewPrometheusExporterInstance(Host string, Port string, Filters map[string]string) (*PrometheusExporter, error) {
	err := checkFiltersMapNotEmpty(Filters)
	if err != nil {
		return nil, err
	}
	return &PrometheusExporter{Host: Host, Port: Port, Filters: Filters}, nil
}

func filtersToString(filters map[string]string) string {
	var filterStrings []string
	for key, value := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s='%s'", key, value))
	}
	return "{" + strings.Join(filterStrings, ", ") + "}"
}

func (p *PrometheusExporter) ExportMetrics() string {
	var wg sync.WaitGroup
	var queryAPIResponse *HTTPResponse
	var metadataAPIResponse *HTTPResponse
	var prometheusQueryAPIResponse PrometheusQueryAPIResponse
	var prometheusMetadataAPIResponse PrometheusMetadataAPIResponse
	wg.Add(2)
	go func(qAPIResp *HTTPResponse) {
		defer wg.Done()
		queryAPIResponse = makeHTTPRequest(p.Host, p.Port, "/api/v1/query", map[string]string{"query": filtersToString(p.Filters)})
	}(queryAPIResponse)

	go func(mAPIResponse *HTTPResponse) {
		defer wg.Done()
		metadataAPIResponse = makeHTTPRequest(p.Host, p.Port, "/api/v1/metadata", map[string]string{})
	}(metadataAPIResponse)

	wg.Wait()
	if queryAPIResponse.Error != nil {
		fmt.Println("Error:", queryAPIResponse.Error)
		return ""
	}
	if metadataAPIResponse.Error != nil {
		fmt.Println("Error:", metadataAPIResponse.Error)
		return ""
	}
	defer queryAPIResponse.Response.Body.Close()
	defer metadataAPIResponse.Response.Body.Close()
	body, err := ioutil.ReadAll(queryAPIResponse.Response.Body)
	if err != nil {
		fmt.Println("Error reading query API response body:", err)
		return ""
	}
	err = json.Unmarshal(body, &prometheusQueryAPIResponse)
	if err != nil {
		fmt.Println("Error decoding JSON received from qeury API:", err)
		return ""
	}
	body, err = ioutil.ReadAll(metadataAPIResponse.Response.Body)
	if err != nil {
		fmt.Println("Error reading query API response body:", err)
		return ""
	}
	err = json.Unmarshal(body, &prometheusMetadataAPIResponse)
	if err != nil {
		fmt.Println("Error decoding JSON received from qeury API:", err)
		return ""
	}
	return exportMetrics(prometheusQueryAPIResponse, prometheusMetadataAPIResponse)
}
