package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/dafvid/timescaleproxy/metric"
)

func (p *Pgdb) createTagsTable(m metric.Metric) *Table {
	t := p.makeTagTable(m)
	qt := "CREATE TABLE IF NOT EXISTS %v(%v)"
	q := fmt.Sprintf(qt, t.Name, t.columnDefs())
	_, err := p.c.Exec(q)
	if err != nil {
		log.Print("Error in db.tags.createTagsTable() ", q, err)
		return nil
	}
	log.Printf("Created tags table for metric '%v'", m.Name)
	return &t
}

func (p Pgdb) makeTagTable(m metric.Metric) Table {
	t := Table{
		Name:   m.Name + "_tags",
		Schema: p.schema,
		Columns: []Column{
			{Name: "id", Type: "serial PRIMARY KEY"},
		},
	}
	t.makeTagColumns(m)

	return t
}

func (t *Table) makeTagColumns(m metric.Metric) []Column {
	result := make([]Column, len(m.Tags))
	for k := range m.Tags {
		c := Column{
			Name: k,
			Type: "text",
		}
		fmt.Println(c)
		t.Columns = append(t.Columns, c)
	}

	return result
}

func (p Pgdb) writeTags(m metric.Metric, t Table) int {
	l := len(t.Columns) - 1
	vf := make([]string, l)
	wf := make([]string, l)
	names := make([]string, l)
	values := make([]interface{}, l)
	var i int8
	for _, c := range t.Columns {
		if c.Name == "id" {
			continue
		}
		names[i] = fmt.Sprintf("\"%v\"", c.Name)
		values[i] = m.Tags[c.Name]
		wf[i] = fmt.Sprintf("\"%v\"=$%v", c.Name, i+1)
		vf[i] = fmt.Sprintf("$%v", i+1)
		i++
	}
	q := fmt.Sprintf("SELECT id FROM %v WHERE %v", t.fullName(), strings.Join(wf, " AND "))
	var tagId int
	row := p.c.QueryRow(q, values...)
	err := row.Scan(&tagId)
	if err == nil {
		return tagId
	}
	//fmt.Println("db.tags.writeTags() ", err)

	q = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v) RETURNING id", t.fullName(), strings.Join(names, ", "), strings.Join(vf, ", "))
	row = p.c.QueryRow(q, values...)
	err = row.Scan(&tagId)
	if err == nil {
		return tagId
	}
	//log.Print("Error in db.tags.writeTags() ", err)
	return 0
}
