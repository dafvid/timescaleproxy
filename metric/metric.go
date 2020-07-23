package metric

import (
	"encoding/json"
	//"fmt"
)

type Metric struct {
	Name      string
	Tags      map[string]string
	Fields    map[string]interface{}
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
