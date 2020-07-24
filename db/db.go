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

func (p *Pgdb) checkTable(m metric.Metric) {
	//fmt.Println("db.checkTable")
	name := m.Name
	//fmt.Println("Name", name)
	_, ok := p.knownTables[name]
	if !ok {
		m.Print()
		var t Table
		if p.Exists(name) {
			t = p.ReflectTable(name)
		} else {
			t = p.CreateTable(m)
		}

		p.knownTables[name] = t
	}
}

func (p Pgdb) Close() error {
	p.connlock.Lock()
	defer p.connlock.Unlock()

	return p.c.Close()
}

func (p *Pgdb) Write(m metric.Metric) {
	//fmt.Println("db.Write()")
	p.checkConn()
	p.checkTable(m)
}
