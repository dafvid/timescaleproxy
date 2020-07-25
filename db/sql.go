package db

import (
	"fmt"
	"log"
	"strings"
	"time"

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

func (t Table) fullName() string {
	return t.Schema + "." + t.Name
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

func (p Pgdb) makeTable(m metric.Metric) Table {
	t := Table{
		Name:   m.Name,
		Schema: p.schema,
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

func (t Table) columnDef() string {
	coldefs := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		coldefs[i] = c.Name + " " + c.Type
	}
	return strings.Join(coldefs, ", ")
}

func (p *Pgdb) createTable(m metric.Metric) *Table {

	t := p.makeTable(m)
	qt := "CREATE TABLE IF NOT EXISTS %v(%v); SELECT create_hypertable('%v','time',chunk_time_interval := '1 week'::interval,if_not_exists := true);"
	q := fmt.Sprintf(qt, t.Name, t.columnDef(), t.Schema+"."+t.Name)
	fmt.Println(q)
	_, err := p.c.Exec(q)
	if err != nil {
		log.Print(err)
		return nil
	}
	log.Printf("Created table for metric '%v'", m.Name)
	return &t
}

func (p *Pgdb) reflectTable(name string) *Table {
	qt := "SELECT column_name, udt_name FROM information_schema.columns WHERE table_schema = '%v' and table_name = '%v'"
	q := fmt.Sprintf(qt, p.schema, name)
	rows, err := p.c.Query(q)
	if err != nil {
		log.Printf("Can't reflect table '%v' (%v)", name, err)
		return nil
	}
	defer rows.Close()
	t := Table{
		Name:    name,
		Schema:  p.schema,
		Columns: []Column{},
	}
	var cname, ctype string
	for rows.Next() {
		err := rows.Scan(&cname, &ctype)
		if err != nil {
			log.Printf("No columns for table '%v''", name)
			return nil
		}
		t.Columns = append(t.Columns, Column{
			Name: cname,
			Type: ctype,
		})
	}
	return &t
}

func (p *Pgdb) exists(name string) bool {
	rows, err := p.c.Query("SELECT 1 FROM pg_tables WHERE tablename = $1 AND schemaname = $2", name, p.schema)
	if err != nil {
		log.Print(err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		return true
	}

	return true
}

func (p Pgdb) write(m metric.Metric) {
	t := p.knownTables[m.Name]
	l := len(t.Columns)
	names := make([]string, l)
	vf := make([]string, l)
	values := make([]interface{}, l)
	for i, c := range t.Columns {
		names[i] = c.Name
		if c.Name == "time" {
			t := time.Unix(m.Timestamp, 0).UTC()
			values[i] = t
		} else {
			values[i] = m.Fields[c.Name]
		}
		vf[i] = fmt.Sprintf("$%v", i+1)
	}
	q := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", t.fullName(), strings.Join(names, ", "), strings.Join(vf, ", "))
	_, err := p.c.Exec(q, values...)
	if err != nil {
		log.Print(err)
	}
}
