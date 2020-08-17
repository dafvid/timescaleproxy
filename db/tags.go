package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"

	"github.com/dafvid/timescaleproxy/log"
	"github.com/dafvid/timescaleproxy/metric"
)

func (p *Pgdb) createTagsTable(ctx context.Context, m metric.Metric) *Table {
	t := p.makeTagTable(m)
	qt := "CREATE TABLE %v(%v)"
	q := fmt.Sprintf(qt, t.Name, t.columnDefs())
	_, err := p.c.Exec(ctx, q)
	if err != nil {
		log.Errorf("Error in db.tags.createTagsTable() '%v' (%v)", q, err)
		return nil
	}
	log.Infof("Created tags table for metric '%v'", m.Name)
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

// Create Table.Columns from Metric.Tags
func (t *Table) makeTagColumns(m metric.Metric) []Column {
	result := make([]Column, len(m.Tags))
	for k := range m.Tags {
		c := Column{
			Name: k,
			Type: "text",
		}
		t.Columns = append(t.Columns, c)
	}

	return result
}

func (p Pgdb) writeTags(ctx context.Context, m metric.Metric, t Table) int {
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
	// See if the tag combination already exist
	q := fmt.Sprintf("SELECT id FROM %v WHERE %v", t.FullName(), strings.Join(wf, " AND "))
	var tagId int
	row := p.c.QueryRow(ctx, q, values...)
	err := row.Scan(&tagId)
	if err == nil {
		return tagId
	} else if err != pgx.ErrNoRows {
		// Print errors besides ErrNoRows
		log.Error(fmt.Sprintf("db.tags.writeTags(): %v (%v)", q, err))
		return 0
	}

	// Otherwise create it
	q = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v) RETURNING id", t.FullName(), strings.Join(names, ", "), strings.Join(vf, ", "))
	row = p.c.QueryRow(ctx, q, values...)
	err = row.Scan(&tagId)
	if err == nil {
		return tagId
	}
	log.Error("Error in db.tags.writeTags() ", err)
	return 0
}
