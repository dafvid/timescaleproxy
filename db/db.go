package db

import (
	//"fmt"
	"sync"

	"github.com/jackc/pgx"

	"github.com/dafvid/timescaleproxy/config"
)

type Pgdb struct {
	connconf pgx.ConnConfig
	c        *pgx.Conn
	connlock sync.Mutex
}

func NewPgdb(conf config.DbConfig) Pgdb {
	p := Pgdb{}
	connconf := pgx.ConnConfig{}
	connconf.Host = conf.Host
	connconf.Database = conf.Database
	connconf.User = conf.User
	connconf.Password = conf.Password

	p.connconf = connconf
	return p
}

func (p Pgdb) connect() error {
	p.connlock.Lock()
	defer p.connlock.Unlock()
	c, err := pgx.Connect(p.connconf)
	if err != nil {
		return err
	}
	p.c = c
	return nil
}

func (p Pgdb) checkConn() {
	p.connlock.Lock()
	defer p.connlock.Unlock()
	if p.c == nil || !p.c.IsAlive() {
		p.connect()
	}
}

func (p Pgdb) Close() error {
	p.connlock.Lock()
	defer p.connlock.Unlock()

	return p.c.Close()
}
