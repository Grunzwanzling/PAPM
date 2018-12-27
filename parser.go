package main

import (
	"flag"
	"gopkg.in/gcfg.v1"
	"os"
	"strings"
)

type Config struct {
	Socket         string
	Db             string
	UseAutoCorrect bool
	Wordlist       string
}

func readFlags() Config {

	wd, _ := os.Getwd()
	var cfg Config
	var cfg_file string
	flag.StringVar(&cfg_file, "config", "./config", "a filepath")
	if cfg_file == "" {
		cfg_file = "./config"
	}
	cfg_file = strings.Replace(cfg_file, "./", wd, -1)
	err := gcfg.ReadFileInto(&cfg, cfg_file)
	if err != nil {

		println("Config parsing error: ", err.Error())
	}
	flag.StringVar(&cfg.Socket, "socket", "./socket", "a filepath")
	flag.StringVar(&cfg.Db, "db", "./db.kdbx", "a kdbx file")
	flag.BoolVar(&cfg.UseAutoCorrect, "use_auto_correct", false, "Wether to auto-correct the password using the wordlist")
	flag.StringVar(&cfg.Wordlist, "wordlist", "./wordlist", "The wordlist to use for auto-correct")

	flag.Parse()

	cfg.Socket = strings.Replace(cfg.Socket, "./", wd, -1)
	return cfg
}
