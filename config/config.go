// From https://paulgorman.org/technical/blog/20171113164018.html

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Configuration struct {
	Db struct {
		Host     string
		Database string
		User     string
		Password string
	}
	Listen struct {
		Address string
		Port    string
	}
}

func Read() (*Configuration, error) {
	fpath := os.Getenv("TIMESCALEPROXY_CONFPATH")
	if fpath == "" {
		return nil, errors.New("Set env TIMESCALEPROXY_CONFPATH")
	}
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
	if err != nil{
		return err
	}
	fmt.Println(fpath)
	fpath = path.Join(fpath, "timescaleproxy.conf.sample")
	fmt.Println(fpath)
	file, err := json.MarshalIndent(Configuration{}, "", "\t")
        if err != nil{
                return err
        }
	fmt.Println(string(file))
	err = ioutil.WriteFile(fpath, file, 0600)
	fmt.Println("Created sample file", fpath)
        if err != nil{
                return err
        }

	return nil
}