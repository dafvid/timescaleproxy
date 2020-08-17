// From https://paulgorman.org/technical/blog/20171113164018.html

package config

import (
	"encoding/json"
	"fmt"
	l "log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/dafvid/timescaleproxy/log"
)

type DbConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	Schema   string
	MaxConns int32
	MinConns int32
}

type Configuration struct {
	Db     DbConfig
	Listen struct {
		Address string
		Port    string
	}
	TimestampUnit     string
	Duration          time.Duration
	DefaultDropPolicy string
	LogLevel          string
}

func Read(cpath string) Configuration {
	fpath := os.Getenv("TIMESCALEPROXY_CONFPATH")
	if fpath == "" {
		fpath = cpath
	}

	if fpath == "" {
		l.Fatal("Set env TIMESCALEPROXY_CONFPATH or use -c flag")
	}
	fpath, err := filepath.Abs(fpath)
	if err != nil {
		l.Fatal(fmt.Sprintf("Config file path error '%v' (%v)", fpath, err))
	}
	file, err := os.Open(fpath)
	if err != nil {
		l.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	c := Configuration{}
	err = decoder.Decode(&c)
	if err != nil {
		l.Fatal(fmt.Sprintf("Can't decode config (%v)", err))
	}

	var lt log.LevelType

	switch c.LogLevel {
	case "off":
		lt = log.Off
	case "error":
		lt = log.ErrorLevel
	case "info":
		lt = log.InfoLevel
	case "debug":
		lt = log.DebugLevel
	}

	log.LogLevel = lt

	var d time.Duration

	if c.TimestampUnit == "" {
		log.Debug("Config 'TimestampUnit' is empty")
		d = time.Second
	} else {
		var err error
		// Stolen from https://github.com/influxdata/telegraf/blob/master/config/config.go
		d, err = time.ParseDuration(c.TimestampUnit)
		if err != nil {
			l.Fatalf("Wrong timestamp unit '%v'. See JSON output docs", c.TimestampUnit)
		}
		// now that we have a duration, truncate it to the nearest
		// power of ten (just in case)
		nearest_exponent := int64(math.Log10(float64(d.Nanoseconds())))
		new_nanoseconds := int64(math.Pow(10.0, float64(nearest_exponent)))
		d = time.Duration(new_nanoseconds)
		// STEAL END
	}
	c.Duration = d
	log.Debugf("Setting timestamp unit to '%v'", d)

	return c
}

func Print() {
	json, err := json.MarshalIndent(Configuration{Db: DbConfig{Schema: "public"}}, "", "\t")
	if err != nil {
		l.Fatal("Can't print Config")
	}
	fmt.Println(string(json))
}
