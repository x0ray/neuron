package main

import (
	l "log"
	"os"

	"github.com/x0ray/neuron"
)

var pgm string

const ver = "0.0.1"

func main() {
	pgm = os.Args[0]
	l.Printf("%s version: %s Started", pgm, ver)
	n := neuron.New()
	l.Printf(n.Status())
	l.Printf("%s Ended", pgm)
}
