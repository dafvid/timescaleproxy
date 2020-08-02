package main

import (
	"flag"
	"fmt"
	l "log"
	"net/http"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/db"
	"github.com/dafvid/timescaleproxy/log"
	"github.com/dafvid/timescaleproxy/metric"
)

var p db.Pgdb

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("\nindex()")
	if r.Body != nil {
		metrics, err := metric.Parse(r.Body)
		if err != nil {
			log.Info("Can't parse JSON ", err)
			http.Error(w, "Can't parse JSON", 400)
		}

		if p.CheckConn() {
			for _, m := range metrics {
				p.Write(m)
			}
			//log.Printf("Wrote %v metrics to DB", len(metrics))
		} else {
			log.Info("Can't connect to backend")
			http.Error(w, "Can't connect to backend", 503)
		}
	}
}

func main() {
	log.Loglevel = log.DebugLevel

	showHelp := flag.Bool("h", false, "Show usage")
	confPath := flag.String("c", "", "Path to config file")
	writeConf := flag.Bool("writeconf", false, "Creates an empty sample conf file")

	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		return
	}
	if *writeConf {
		err := config.Write()
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	conf, err := config.Read(*confPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	p = db.NewPgdb(conf.Db)
	http.HandleFunc("/", index)
	listenStr := conf.Listen.Address + ":" + conf.Listen.Port
	log.Info(fmt.Sprint("Starting server ", listenStr))
	l.Fatal(http.ListenAndServe(listenStr, nil))
}
