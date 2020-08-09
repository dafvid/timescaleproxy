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
}

func NewPgdb(conf config.DbConfig) *Pgdb {
	p := Pgdb{}

	poolconf, err := pgxpool.ParseConfig("")
	if err != nil {
		l.Fatal("Cannot parse config", err)
	}
	poolconf.MinConns = conf.MinConns
	poolconf.MaxConns = conf.MaxConns
	poolconf.ConnConfig.Host = conf.Host
	if conf.Port == 0 {
		conf.Port = 5432
	}
	poolconf.ConnConfig.Port = conf.Port
	poolconf.ConnConfig.Database = conf.Database
	poolconf.ConnConfig.User = conf.User
	poolconf.ConnConfig.Password = conf.Password

	log.Print(*poolconf)

	log.Info(fmt.Sprintf("Connecting to Postgresql (%v:%v)", conf.Host, conf.Port))
	c, err := pgxpool.ConnectConfig(context.Background(), poolconf)
	if err != nil {
		log.Print("db.connect() ", err)
		return nil
	}
	p.c = c
	p.knownTables = make(map[string]Tables)
	p.schema = conf.Schema
	return &p
}

func (p *Pgdb) checkTables(ctx context.Context, m metric.Metric) bool {
	//fmt.Println("db.checkTable")
	name := m.Name
	//fmt.Println("Name", name)
	_, ok := p.knownTables[name]
	if !ok {
		//m.Print()
		var t *Tables

		if p.exists(ctx, name) {
			t = p.reflectTables(ctx, name)
		} else {
			t = p.createTables(ctx, m)
		}

		if t == nil {
			return false
		}

		p.knownTables[name] = *t
	}

	return true
}

func (p *Pgdb) Write(ctx context.Context, m metric.Metric) bool {
	//fmt.Println("db.Write()")
	ok := p.checkTables(ctx, m)
	if !ok {
		//log.Printf("Can't create table for metric '%v'", m.Name)
		return false
	} else {
		return p.write(ctx, m)
	}
}
