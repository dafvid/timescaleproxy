// From https://paulgorman.org/technical/blog/20171113164018.html

package config

import (
	"encoding/json"
	"fmt"
	l "log"
	"os"
	"path/filepath"
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

	return c
}

func Print() {
	json, err := json.MarshalIndent(Configuration{Db: DbConfig{Schema: "public"}}, "", "\t")
	if err != nil {
		l.Fatal("Can't print Config")
	}
	fmt.Println(string(json))
}
