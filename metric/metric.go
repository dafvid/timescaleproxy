package metric

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	//"github.com/dafvid/timescaleproxy/log"
)

type Tags map[string]string
type Fields map[string]interface{}

type Metric struct {
	Name      string
	Tags      Tags
	Fields    Fields
	Timestamp int64
}

type Metrics struct {
	Metrics []Metric
}

func (m Metric) MissingValues() bool {
	return m.Name == "" || len(m.Fields) == 0
}

func Parse(data io.Reader) ([]Metric, error) {
	peekReader := bufio.NewReader(data)
	peek, _ := peekReader.Peek(9) // peek first 9 bytes of content
	isBatch := string(peek[2:]) == "metrics"
	decoder := json.NewDecoder(peekReader)

	if isBatch {
		var metrics Metrics
		err := decoder.Decode(&metrics)
		if err != nil {
			return []Metric{}, err
		}
		return metrics.Metrics, nil
	} else {
		var metric Metric
		err := decoder.Decode(&metric)
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
