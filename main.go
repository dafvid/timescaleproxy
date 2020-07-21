package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dafvid/timescaleproxy/config"
	"io"
	"log"
	"net/http"
)

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
	//body, _ := r.GetBody()
	fmt.Println("Request incoming!")
	if r.Body != nil {
		var b bytes.Buffer
		var dest io.Writer = &b
		_, _ = io.Copy(dest, r.Body)
		//fmt.Println(string(b.Bytes()))
		var result map[string]interface{}
		json.Unmarshal(b.Bytes(), &result)
		_, isBatch := result["metrics"]
		if isBatch {
			for _, m := range result["metrics"].([]interface{}) {
				handleMetric(m.(map[string]interface{}))
			}
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
	http.HandleFunc("/", index)
	listenStr := conf.Listen.Address + ":" + conf.Listen.Port
	log.Print("Starting server ", listenStr)
	log.Fatal(http.ListenAndServe(listenStr, nil))
}
