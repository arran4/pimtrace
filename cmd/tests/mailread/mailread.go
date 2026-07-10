package main

import (
	"flag"
	"log"
	"pimtrace/dataformats"
	"pimtrace/dataformats/maildata"
	"pimtrace/fsys"
)

/**
For bug finding with mail.
*/

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Printf("No files specified")
	}
	for _, fn := range flag.Args() {
		log.Printf("Reading %s", fn)
		nm, err := dataformats.ReadFile(fsys.OSFS{}, "mailfile", fn, maildata.ReadMailStream)
		if err != nil {
			log.Panicf("Read error: %s", err)
		}
		log.Print("Read", len(nm))
	}
}
