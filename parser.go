package main

import (
	"flag"
	"gopkg.in/gcfg.v1"
	"os"
	"strings"
)

type Config struct {
	Socket         string
	DbList         []string
	UseAutoCorrect bool
	Wordlist       string
}

func readFlags() Config {

	var cfg Config
	cfg_file := flag.String("config", "./config", "a filepath")
	err := gcfg.ReadFileInto(&cfg, *cfg_file)
	if err != nil {

		println("Config parsing error: ", err.Error())
	}
	flag.StringVar(&cfg.Socket, "socket", "~/socket", "a filepath")
	flag.BoolVar(&cfg.UseAutoCorrect, "use_auto_correct", false, "Wether to auto-correct the password using the wordlist")
	flag.StringVar(&cfg.Wordlist, "wordlist", "./wordlist", "The wordlist to use for auto-correct")

	flag.Parse()

	wd, _ := os.Getwd()
	cfg.Socket = strings.Replace(cfg.Socket, "./", wd, -1)
	return cfg
}
