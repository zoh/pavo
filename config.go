package main

import (
	"os"

	"code.google.com/p/getopt"
)

type Config struct {
	Root   string
	Host   string
	CORS   bool
	Static bool
}

func ParseArgs() *Config {
	help := getopt.BoolLong("help", 'h', "Show help")
	root := getopt.StringLong("root", 'r', "./dummy/root_storage", "Root of storage")
	host := getopt.StringLong("host", 's', "localhost:9073", "host:port for pavo server")
	cors := getopt.BoolLong("cors", 'c', "Enable CORS headers")
	static := getopt.BoolLong("static", 'l', "Serve uploaded files")
	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(1)
	}

	return &Config{
		Root:   *root,
		Host:   *host,
		CORS:   *cors,
		Static: *static,
	}
}
