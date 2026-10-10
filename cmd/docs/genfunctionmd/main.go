package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Create("functions.md")
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	_, err = f.WriteString(generateMarkdown())
	if err != nil {
		log.Panicln(err)
	}
}
