package db

import (
	"context"
	"fmt"
	l "log"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/log"
	"github.com/dafvid/timescaleproxy/metric"
)

type Pgdb struct {
	c           *pgxpool.Pool
	knownTables map[string]Tables
	schema      string
	config      config.Configuration
}

func NewPgdb(conf config.Configuration) *Pgdb {
	p := Pgdb{
		config: conf,
	}

	poolconf, err := pgxpool.ParseConfig("")
	if err != nil {
		l.Fatal("Cannot parse pool config", err)
	}
	if conf.Db.MinConns != 0 {
		poolconf.MinConns = conf.Db.MinConns
	}
	if conf.Db.MaxConns != 0 {
		poolconf.MaxConns = conf.Db.MaxConns
	}
	poolconf.ConnConfig.Host = conf.Db.Host
	if conf.Db.Port != 0 {
		poolconf.ConnConfig.Port = conf.Db.Port
	}
	poolconf.ConnConfig.Database = conf.Db.Database
	poolconf.ConnConfig.User = conf.Db.User
	poolconf.ConnConfig.Password = conf.Db.Password

	log.Info(fmt.Sprintf("Connecting to Postgresql (%v:%v) as %v", conf.Db.Host, conf.Db.Port, conf.Db.User))
	c, err := pgxpool.ConnectConfig(context.Background(), poolconf)
	if err != nil {
		log.Error("db.connect() ", err)
		return nil
	}
	p.c = c
	p.knownTables = make(map[string]Tables)
	p.schema = conf.Db.Schema
	return &p
}

func (p *Pgdb) checkTables(ctx context.Context, m metric.Metric) bool {
	name := m.Name
	tagsName := name + "_tags"
	_, ok := p.knownTables[name]
	if !ok {
		t := Tables{}

		var tt *Table
		if p.tableExists(ctx, tagsName) {
			tt = p.reflectTable(ctx, tagsName)
		} else {
			tt = p.createTagsTable(ctx, m)
		}
		if tt == nil {
			return false
		}
		t.TagsTable = *tt

		var dt *Table
		if p.tableExists(ctx, name) {
			dt = p.reflectTable(ctx, name)
		} else {
			dt = p.createDataTable(ctx, m)
		}
		if dt == nil {
			return false
		}
		t.DataTable = *dt

		p.knownTables[name] = t
	}

	return true
}

func (p *Pgdb) Write(ctx context.Context, m metric.Metric) bool {
	//fmt.Println("db.Write()")
	ok := p.checkTables(ctx, m)
	if !ok {
		return false
	} else {
		return p.write(ctx, m)
	}
}
