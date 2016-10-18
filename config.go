package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	UID           string
	AuthToken     string
	APIToken      string
	HelpRequested bool
	IgnoreDonuts  bool
	ConfigFile    string
	LogDir        string
}

func (this *Config) String() string {
	banner := []string{
		"\r\n\tLevel User ID  => ", this.UID,
		"\r\n\tLevel Auth Token => ", this.AuthToken,
		"\r\n\tLevel API Token => ", this.APIToken,
		"\r\n\tConfiguration File=> ", this.ConfigFile,
		"\r\n\tLog Directory => ", this.LogDir,
		"\r\n\tBase URL => ", ConfigurationFileMap.BaseURL,
		"\r\n\tIgnore Donuts => ", strconv.FormatBool(this.IgnoreDonuts),
	}
	return strings.Join(banner, " ")
}

var (
	Configuration = &Config{
		UID:           "",
		AuthToken:     "",
		APIToken:      "",
		HelpRequested: false,
		IgnoreDonuts:  false,
		ConfigFile:    os.Getenv("GOPATH") + "/src/github.com/javand/go_level_money/configuration.json",
		LogDir:        "log",
	}
)

type ConfigurationFile struct {
	BaseURL string
}

var ConfigurationFileMap ConfigurationFile

func ParseCLIArgs() {
	flag.Parse()
	if Configuration.HelpRequested {
		flag.Usage()
		os.Exit(0)
	}
	log.Print("Using configuration settings:", Configuration)
}

func LoadConfigurationFile() {
	file, err := ioutil.ReadFile(Configuration.ConfigFile)
	if err != nil {
		log.Panicln("File error: ", err)
	}

	err = json.Unmarshal(file, &ConfigurationFileMap)
	if err != nil {
		log.Panicln("JSON decoding error: ", err)
	}
}
