package main

import (
	"flag"
	"fmt"
	"html"
	l "log"
	"net/http"
	"strings"
	"time"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/db"
	"github.com/dafvid/timescaleproxy/log"
	"github.com/dafvid/timescaleproxy/metric"
)

var p db.Pgdb

func index(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	ct := r.Header.Get("content-type")
	if r.Body == http.NoBody || ct == "" {
		fmt.Fprintln(w,
			"Blip! Blop! Blip! This here is TimescaleDB HTTP Proxy v0.1! "+html.EscapeString("<<END OF TRANSMISSION"))
	} else {
		cta := strings.Split(ct, ";")
		if len(cta) > 1 {
			ct = cta[0]
		}
		switch ct {
		case "application/json":
			metrics, err := metric.Parse(r.Body)
			if err != nil {
				log.Error("main.index(): Can't parse JSON ", err)
				http.Error(w, "Can't parse JSON", 400)
				return
			}

			for _, m := range metrics {
				if !p.Write(r.Context(), m) {
					http.Error(w, fmt.Sprintf("Can't write metrics '%v'", m.Name), 503)
					return
				}
			}
			log.Debug(fmt.Sprintf("Wrote %3d metrics to DB in %5.3f s", len(metrics), time.Since(t).Seconds()))
		default:
			http.Error(w, fmt.Sprintf("Unsupported Content-type '%v'", r.Header.Get("content-type")), 400)
			return
		}
	}
}

func main() {
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
	log.PrintLevel()
	tp := db.NewPgdb(conf)
	if tp == nil {
		l.Fatal("Cannot create backend")
	}
	p = *tp
	http.HandleFunc("/", index)
	listenStr := conf.Listen.Address + ":" + conf.Listen.Port
	log.Info(fmt.Sprint("Starting server ", listenStr))
	l.Fatal(http.ListenAndServe(listenStr, nil))
}
