// From https://paulgorman.org/technical/blog/20171113164018.html

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type DbConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	Schema   string
}

type Configuration struct {
	Db     DbConfig
	Listen struct {
		Address string
		Port    string
	}
	TimestampUnit string
}

func Read(cpath string) (*Configuration, error) {
	fpath := os.Getenv("TIMESCALEPROXY_CONFPATH")
	if fpath == "" {
		fpath = cpath
	}

	if fpath == "" {
		return nil, errors.New("Set env TIMESCALEPROXY_CONFPATH or use -c flag")
	}
	fpath, err := filepath.Abs(fpath)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		return nil, err
	}

	return &Config, nil
}

func Write() error {
	fpath, err := os.Getwd()
	if err != nil {
		return err
	}
	fpath = path.Join(fpath, "timescaleproxy.conf.sample")
	file, err := json.MarshalIndent(Configuration{Db: DbConfig{Schema: "public"}}, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fpath, file, 0600)
	log.Print("Created sample file", fpath)
	if err != nil {
		return err
	}

	return nil
}
