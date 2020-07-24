package db

import (
	"fmt"
	"log"
	"strings"

	//"github.com/jackc/pgx"

	"github.com/dafvid/timescaleproxy/metric"
)

type Column struct {
	Name string
	Type string
}

type Table struct {
	Name    string
	Schema  string
	Columns []Column
}

func getDataType(v interface{}) string {
	switch v.(type) {
	case bool:
		return "boolean"
	case string:
		return "text"
	case float64, float32:
		return "float8"
	default:
		log.Printf("Unknown type $T($v)", v, v)
		return "text"
	}
}

func makeTable(m metric.Metric) Table {
	t := Table{
		Name:   m.Name,
		Schema: "public",
		Columns: []Column{
			{Name: "time", Type: "timestamp"},
		},
	}
	t.makeColumns(m)

	return t
}

func (t *Table) makeColumns(m metric.Metric) []Column {
	result := make([]Column, len(m.Fields))
	for k, v := range m.Fields {
		t.Columns = append(t.Columns, Column{
			Name: k,
			Type: getDataType(v),
		})
	}

	return result
}

func (t Table) ColumnDef() string {
	coldefs := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		coldefs[i] = c.Name + " " + c.Type
	}
	return strings.Join(coldefs, ", ")
}

func (p *Pgdb) CreateTable(m metric.Metric) Table {
	t := makeTable(m)
	qf := "CREATE TABLE IF NOT EXISTS %v(%v); SELECT create_hypertable('%v','time',chunk_time_interval := '1 week'::interval,if_not_exists := true);"
	q := fmt.Sprintf(qf, t.Name, t.ColumnDef(), t.Schema+"."+t.Name)
	fmt.Println(q)
	/*ct, err := p.c.Exec(q)
	  if err != nil {
	      log.Print(err)
	  }*/
	return t
}

func (p *Pgdb) ReflectTable(name string) Table {
	t := Table{}
	return t
}

func (p *Pgdb) Exists(name string) bool {
	ct, err := p.c.Exec("SELECT tablename FROM pg_tables WHERE tablename = $1 AND schemaname = 'public'", name)
	if err != nil {
		log.Print(err)
		return false
	}

	if ct.RowsAffected() == 1 {
		return true
	}

	return false
}
