package main

import (
	l "log"
	"os"

	u "github.com/x0ray/neuron"
)

var pgm string

const ver = "0.0.1"

func main() {
	pgm = os.Args[0]
	l.Printf("%s version: %s Started", pgm, ver)

	// make neurons and set their initial thresholds
	in := u.New(5)
	plain := u.NewNeurons(10, 3)
	out := u.New(22)

	// connect neurons
	in.LinkOneToMany(plain)
	out.LinkManyToOne(plain)

	// start all neurons
	u.StartAllNeurons()

	// fire a neuron
	in.Fire(7)

	// print status of all neurons
	u.Status()

	// stop all neurons
	u.StopAllNeurons()

	l.Printf("%s Ended", pgm)
}
