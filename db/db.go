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
	knownTables map[string]Tables
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
	p.knownTables = make(map[string]Tables)
	p.schema = conf.Schema
	return p
}

func (p *Pgdb) IsAlive() bool {
	if !p.c.IsAlive() {
		return false
	}
	row := p.c.QueryRow("SELECT 1")
	var i int
	err := row.Scan(&i)
	if err != nil {
		log.Print("Connection wasn't alive really...")
		return false
	}
	return true
}

func (p *Pgdb) connect() *pgx.Conn {
	log.Printf("Connecting to Postgresql (%v:%v)", p.connconf.Host, p.connconf.Port)
	c, err := pgx.Connect(p.connconf)
	if err != nil {
		log.Print("db.connect() ", err)
	}
	return c
}

func (p *Pgdb) CheckConn() bool {
	//fmt.Println("db.checkConn()")
	p.connlock.Lock()
	defer p.connlock.Unlock()
	if p.c == nil {
		p.c = p.connect()
	} else if !p.IsAlive() {
		p.c.Close()
		p.c = p.connect()
	}
	return p.c != nil
}

func (p *Pgdb) checkTables(m metric.Metric) bool {
	//fmt.Println("db.checkTable")
	name := m.Name
	//fmt.Println("Name", name)
	_, ok := p.knownTables[name]
	if !ok {
		//m.Print()
		var t *Tables

		if p.exists(name) {
			t = p.reflectTables(name)
		} else {
			t = p.createTables(m)
		}

		if t == nil {
			return false
		}

		p.knownTables[name] = *t
	}

	return true
}

func (p *Pgdb) Write(m metric.Metric) {
	//fmt.Println("db.Write()")
	ok := p.checkTables(m)
	if !ok {
		log.Printf("No table for metric '%v'", m.Name)
	} else {
		p.write(m)
	}
}
