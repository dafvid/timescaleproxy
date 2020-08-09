package main

import (
	"flag"
	"fmt"
	l "log"
	"net/http"
	"time"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/db"
	"github.com/dafvid/timescaleproxy/log"
	"github.com/dafvid/timescaleproxy/metric"
)

var p db.Pgdb

func index(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	//fmt.Println("\nindex()")
	if r.Body != nil {
		metrics, err := metric.Parse(r.Body)
		if err != nil {
			log.Info("main.index(): Can't parse JSON ", err)
			http.Error(w, "Can't parse JSON", 400)
		}

		for _, m := range metrics {
			if !p.Write(r.Context(), m) {
				http.Error(w, fmt.Sprintf("Can't write metrics '%v'", m.Name), 503)
			}
		}
		log.Debug(fmt.Sprintf("Wrote %3d metrics to DB in %5.3f s", len(metrics), time.Since(t).Seconds()))
	}
}

func main() {
	log.Loglevel = log.DebugLevel

	showHelp := flag.Bool("h", false, "Show usage")
	confPath := flag.String("c", "", "Path to config file")
	printConf := flag.Bool("printconf", false, "Prints empty sample conf to stdout")

	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		return
	}
	if *printConf {
		config.Print()
		return
	}
	conf := config.Read(*confPath)
	tp := db.NewPgdb(conf.Db)
	if tp == nil {
		l.Fatal("Cannot connect to database")
	}
	p = *tp
	http.HandleFunc("/", index)
	listenStr := conf.Listen.Address + ":" + conf.Listen.Port
	log.Info(fmt.Sprint("Starting server ", listenStr))
	l.Fatal(http.ListenAndServe(listenStr, nil))
}
