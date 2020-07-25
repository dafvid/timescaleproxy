package db

import (
	"log"
	"sync"

	"github.com/jackc/pgx"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/metric"
)

type Pgdb struct {
	connconf    pgx.ConnConfig
	c           *pgx.Conn
	connlock    sync.Mutex
	knownTables map[string]Table
	schema      string
}

func NewPgdb(conf config.DbConfig) Pgdb {
	p := Pgdb{}
	connconf := pgx.ConnConfig{}
	connconf.Host = conf.Host
	connconf.Port = conf.Port
	connconf.Database = conf.Database
	connconf.User = conf.User
	connconf.Password = conf.Password

	p.connconf = connconf
	p.knownTables = make(map[string]Table)
	p.schema = conf.Schema
	return p
}

func (p *Pgdb) checkConn() error {
	//fmt.Println("db.checkConn()")
	p.connlock.Lock()
	defer p.connlock.Unlock()
	if p.c == nil || !p.c.IsAlive() {
		log.Printf("Connecting to DB (%v:%v)...", p.connconf.Host, p.connconf.Port)
		c, err := pgx.Connect(p.connconf)
		if err != nil {
			log.Print(err)
		}
		p.c = c
		log.Printf("...done!")
	}
	return nil
}

func (p *Pgdb) checkTable(m metric.Metric) *Table {
	//fmt.Println("db.checkTable")
	name := m.Name
	//fmt.Println("Name", name)
	t, ok := p.knownTables[name]
	if !ok {
		//m.Print()
		if p.exists(name) {
			t = *p.reflectTable(name)
		} else {
			t = *p.createTable(m)
		}

		p.knownTables[name] = t
	}

	return &t
}

func (p Pgdb) Close() error {
	p.connlock.Lock()
	defer p.connlock.Unlock()

	return p.c.Close()
}

func (p *Pgdb) Write(m metric.Metric) {
	//fmt.Println("db.Write()")
	p.checkConn()
	t := p.checkTable(m)
	if t == nil {
		log.Print("No table for metric '%v'", m.Name)
	} else {
		p.write(m)
	}
}
