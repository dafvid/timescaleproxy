package metric

import (
	"encoding/json"
	"fmt"
)

type Metric struct {
	name      string
	tags      map[string]string
	fields    map[string]interface{}
	timestamp int
}

type Metrics struct {
	metrics []Metric
}

func Parse(data []byte) ([]Metric, error) {
	fmt.Println(string(data))
	// peek into first bytes of data
	isBatch := string(data[2:9]) == "metrics"
	if isBatch {
		fmt.Println("It's a batch!")
		var result Metrics
		err := json.Unmarshal(data, &result)
		if err != nil {
			return []Metric{}, err
		}
		return result.metrics, nil
	} else {
		var result Metric
		err := json.Unmarshal(data, &result)
		if err != nil {
			return []Metric{}, err
		}
		return []Metric{result}, nil
	}
}
