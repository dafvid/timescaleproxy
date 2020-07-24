package metric

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Tags map[string]string
type Fields map[string]interface{}

type Metric struct {
	Name      string
	Tags      Tags
	Fields    Fields
	Timestamp int
}

type Metrics struct {
	Metrics []Metric
}

func Parse(data []byte) ([]Metric, error) {
	// peek into first bytes of data
	isBatch := string(data[2:9]) == "metrics"
	if isBatch {
		var metrics Metrics
		err := json.Unmarshal(data, &metrics)
		if err != nil {
			return []Metric{}, err
		}
		return metrics.Metrics, nil
	} else {
		var metric Metric
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return []Metric{}, err
		}
		return []Metric{metric}, nil
	}
}

func (m Metric) Print() {
	fmt.Println()
	fmt.Println("METRIC")
	fmt.Println("  name =", m.Name)
	fmt.Println("  ts =", int(m.Timestamp))
	fmt.Println("FIELDS")
	for k, v := range m.Fields {
		fmt.Println(" ", k, "=", v, reflect.TypeOf(v))
	}

	fmt.Println("TAGS")
	for k, v := range m.Tags {
		fmt.Println(" ", k, "=", v)
	}

}
