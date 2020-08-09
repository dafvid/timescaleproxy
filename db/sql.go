package db

import (
	"context"
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

type Tables struct {
	DataTable Table
	TagsTable Table
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

// Create Table from Metric
func (p Pgdb) makeDataTable(m metric.Metric) Table {
	t := Table{
		Name:   m.Name,
		Schema: p.schema,
		Columns: []Column{
			{Name: "time", Type: "timestamp"},
			{Name: "tag_id", Type: "integer REFERENCES " + m.Name + "_tags (id)"},
		},
	}
	t.makeColumns(m)

	return t
}

// Create Columns from Metric.Fields
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

// Create SQL from Table.Columns
func (t Table) columnDefs() string {
	coldefs := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		coldefs[i] = "\"" + c.Name + "\" " + c.Type
	}
	return strings.Join(coldefs, ", ")
}

func (p *Pgdb) createDataTable(ctx context.Context, m metric.Metric) *Table {
	t := p.makeDataTable(m)
	qt := "CREATE TABLE IF NOT EXISTS %v(%v); SELECT create_hypertable('%v','time',chunk_time_interval := '1 week'::interval,if_not_exists := true);"
	q := fmt.Sprintf(qt, t.Name, t.columnDefs(), t.Schema+"."+t.Name)
	_, err := p.c.Exec(ctx, q)
	if err != nil {
		log.Print("Error in db.sql.createTable() ", q, err)
		return nil
	}
	log.Printf("Created table for metric '%v'", m.Name)
	return &t
}

func (p *Pgdb) createTables(ctx context.Context, m metric.Metric) *Tables {
	//log.Printf("db.sql.createTables(%v)", m.Name)
	tt := p.createTagsTable(ctx, m)
	if tt == nil {
		return nil
	}
	t := p.createDataTable(ctx, m)
	if t == nil {
		return nil
	}

	return &Tables{
		TagsTable: *tt,
		DataTable: *t,
	}
}

func (p *Pgdb) reflectTable(ctx context.Context, name string) *Table {
	qt := "SELECT column_name, udt_name FROM information_schema.columns WHERE table_schema = '%v' and table_name = '%v'"
	q := fmt.Sprintf(qt, p.schema, name)
	rows, err := p.c.Query(ctx, q)
	defer rows.Close()
	if err != nil {
		log.Printf("Error in db.sql.reflectTable(): %v (%v)", q, err)
		return nil
	}

	t := Table{
		Name:    name,
		Schema:  p.schema,
		Columns: []Column{},
	}
	var cname, ctype string
	for rows.Next() {
		err := rows.Scan(&cname, &ctype)
		if err != nil {
			log.Printf("Error in db.sql.reflectTable(): No columns for table '%v''", name)
			return nil
		}
		t.Columns = append(t.Columns, Column{
			Name: cname,
			Type: ctype,
		})
	}
	return &t
}

func (p *Pgdb) reflectTables(ctx context.Context, name string) *Tables {
	//log.Printf("db.sql.reflectTables(%v)", name)
	tt := p.reflectTable(ctx, name+"_tags")
	if tt == nil {
		return nil
	}
	dt := p.reflectTable(ctx, name)
	if dt == nil {
		return nil
	}

	return &Tables{
		TagsTable: *tt,
		DataTable: *dt,
	}
}

func (p *Pgdb) tableExists(ctx context.Context, name string) bool {
	//log.Printf("db.sql.exists(%v)", name)
	qt := "SELECT 1 FROM pg_tables WHERE tablename = '%v' AND schemaname = '%v'"
	q := fmt.Sprintf(qt, name, p.schema)
	rows, err := p.c.Query(ctx, q)
	defer rows.Close()
	if err != nil {
		log.Print("Error in db.tableExists() ", q, err)
		return false
	}

	for rows.Next() {
		return true
	}

	return false
}

func (p *Pgdb) exists(ctx context.Context, name string) bool {
	return p.tableExists(ctx, name) && p.tableExists(ctx, name+"_tags")
}

func (p Pgdb) writeData(ctx context.Context, m metric.Metric, t Table, tagId int) bool {
	l := len(t.Columns)
	names := make([]string, l)
	vf := make([]string, l)
	values := make([]interface{}, l)
	for i, c := range t.Columns {
		names[i] = fmt.Sprintf("\"%v\"", c.Name)
		if c.Name == "time" {
			t := time.Unix(m.Timestamp, 0).UTC()
			values[i] = t
		} else if c.Name == "tag_id" {
			values[i] = tagId
		} else {
			values[i] = m.Fields[c.Name]
		}
		vf[i] = fmt.Sprintf("$%v", i+1)
	}
	q := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", t.fullName(), strings.Join(names, ", "), strings.Join(vf, ", "))
	_, err := p.c.Exec(ctx, q, values...)
	if err != nil {
		log.Print("Error in db.writeData()", q, err)
		return false
	}
	return true
}

func (p Pgdb) write(ctx context.Context, m metric.Metric) bool {
	t := p.knownTables[m.Name]
	tagId := p.writeTags(ctx, m, t.TagsTable)
	if tagId == 0 {
		//log.Print("Can't find tag id for ", m.Name)
		//m.Print()
		return false
	}
	return p.writeData(ctx, m, t.DataTable, tagId)
}
