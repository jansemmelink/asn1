package main

import (
	"flag"

	"bitbucket.org/vservices/dark/logger"
	"bitbucket.org/vservices/dark/logger/level"
	"github.com/jansemmelink/asn1/asn1def"
	"github.com/pkg/errors"
)

var log = logger.New()

func main() {
	fileFlag := flag.String("file", "", "ASN.1 definition file to load")
	debugFlag := flag.Bool("d", false, "DEBUG mode")
	flag.Parse()
	if *debugFlag {
		logger.Top().WithLevel(level.Debug)
	}

	err := asn1def.New().LoadFile(*fileFlag)
	if err != nil {
		panic(errors.Wrap(err, "Failed to load file "+*fileFlag))
	}
}
