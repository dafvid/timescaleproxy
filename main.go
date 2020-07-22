package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dafvid/timescaleproxy/config"
	"github.com/dafvid/timescaleproxy/db"
	"github.com/dafvid/timescaleproxy/metric"
)

var p db.Pgdb

func handleMetric(m map[string]interface{}) {
	fmt.Println()
	fmt.Println("METRIC")
	fmt.Println("  name =", m["name"])
	fmt.Println("  ts =", int(m["timestamp"].(float64)))
	fmt.Println("FIELDS")
	for k, v := range m["fields"].(map[string]interface{}) {
		fmt.Println(" ", k, "=", v)
	}

	fmt.Println("TAGS")
	for k, v := range m["tags"].(map[string]interface{}) {
		fmt.Println(" ", k, "=", v)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		var b bytes.Buffer
		var dest io.Writer = &b
		_, _ = io.Copy(dest, r.Body)
		metrics, err := metric.Parse(b.Bytes())
		if err != nil {
			log.Print(err)
		}
		for _, m := range metrics {
			fmt.Println(m)
		}
	}
}

func main() {
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
	log.Print("Starting server ", listenStr)
	log.Fatal(http.ListenAndServe(listenStr, nil))
}
