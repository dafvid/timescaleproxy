package db

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/jackc/pgx"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/metric"
	"github.com/dafvid/timescaleproxy/util"
)

type Pgdb struct {
	connconf pgx.ConnConfig
	c        *pgx.Conn
	connlock sync.Mutex
}

var tableChecked []string

func NewPgdb(conf config.DbConfig) Pgdb {
	p := Pgdb{}
	connconf := pgx.ConnConfig{}
	connconf.Host = conf.Host
	connconf.Port = conf.Port
	connconf.Database = conf.Database
	connconf.User = conf.User
	connconf.Password = conf.Password

	p.connconf = connconf
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

func printMetric(m metric.Metric) {
	fmt.Println()
	fmt.Println("METRIC")
	fmt.Println("  name =", m.Name)
	fmt.Println("  ts =", int(m.Timestamp))
	fmt.Println("FIELDS")
	for k, v := range m.Fields {
		fmt.Println(" ", k, "=", v, reflect.TypeOf(v))
	}

	fmt.Println("TAGS")
	for k, v := range m.Tags {
		fmt.Println(" ", k, "=", v)
	}

}

func (p Pgdb) checkTable(m metric.Metric) {
	//fmt.Println("db.checkTable")
	name := m.Name
	//fmt.Println("Name", name)
	if !util.InArr(name, tableChecked) {
		// table exists
		// else create table
		printMetric(m)
		tableChecked = append(tableChecked, name)
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
