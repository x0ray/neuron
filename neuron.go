package neuron

import (
	"fmt"
	"log"
	"reflect"
)

type nstate uint8

const (
	stopped nstate = iota
	running
	paused
)

var nstates = [...]string{
	"stopped",
	"running",
	"paused",
}

func (e nstate) String() string {
	if e < stopped || e > paused {
		return "unknown"
	}
	return nstates[e]
}

type cntlMsg struct {
	cmd string
}

var (
	neurons   []*Neuron
	cntlChans []chan cntlMsg
)

// Neuron - simulated neuron, composed of this dats structure and
//   its own Go routine (when in running state)
type Neuron struct {
	cntlChan      chan cntlMsg // input from central neuron controller
	dendriteChans []chan int   // dendrite is input from other neurons
	axonChans     []chan int   // axon is output to other neurons
	connected     bool         // connected to other neurons
	state         nstate       // can be: stopped, running, paused
	sigma         int          // sum of dendrite signals
	threshold     int          // level required for axon output
	lastInChan    chan int
}

// New - creates a new neuron in the stopped state
func New() *Neuron {
	n := new(Neuron)
	n.cntlChan = make(chan cntlMsg)
	cntlChans = append(cntlChans, n.cntlChan)
	n.connected = false
	n.state = stopped
	neurons = append(neurons, n)
	return n
}

// NewNeurons - create many neurons in a slice
func NewNeurons(num int) []*Neuron {
	ns := make([]*Neuron, num)
	for i := 0; i < num; i++ {
		ns[i] = New()
	}
	return ns
}

// Status - print neuron status
func (n *Neuron) Status() string {
	var s string
	s = fmt.Sprintf("\nNeuron: %p\n  State..........: %v\n  Dendrite inputs: %d\n  Axon outputs...: %d\n"+
		"  Connected......: %v\n  Sigma..........: %d\n  Threshold......: %d",
		n, n.state, len(n.dendriteChans), len(n.axonChans), n.connected, n.sigma, n.threshold)
	return s
}

// Link - connect an axom channel of a neuron to a dendrite channel of a neuron
func (n *Neuron) Link(to *Neuron) *Neuron {
	a := make(chan int)
	to.dendriteChans = append(to.dendriteChans, a)
	to.connected = true
	n.axonChans = append(n.axonChans, a)
	n.connected = true
	return to
}

// LinkManyToOne - connect an axom channel of many neurons to dendrite channels of a neuron
func (n *Neuron) LinkManyToOne(many []*Neuron) *Neuron {
	for _, v := range many {
		a := make(chan int)
		v.axonChans = append(v.axonChans, a)
		v.connected = true
		n.dendriteChans = append(n.dendriteChans, a)
		n.connected = true
	}
	return n
}

// LinkOneToMany - connect an axom channel of one neuron to dendrite channels of many neurons
func (n *Neuron) LinkOneToMany(many []*Neuron) *Neuron {
	for _, v := range many {
		a := make(chan int)
		v.connected = true
		n.axonChans = append(n.axonChans, a)
		n.connected = true
		v.dendriteChans = append(v.dendriteChans, a)
	}
	return n
}

// Fire - send an impulse to a specific neuron
func (n *Neuron) Fire(impulse int) *Neuron {
	if len(n.dendriteChans) > 0 {
		fireChan := n.dendriteChans[0]
		fireChan <- impulse
	}
	return n
}

// StartAllNeurons - start all the defined neurons
func StartAllNeurons() {
	for i, v := range neurons {
		v.Run()
		log.Printf("Started neuron%d: %p", i, v)
	}
}

// StopAllNeurons - start all the defined neurons
func StopAllNeurons() {
	for i, v := range cntlChans {
		v <- cntlMsg{cmd: "stop"}
		log.Printf("Stopped neuron%d on cntl chan: %p", i, v)
	}
}

// Run - start all neuron channels executing
func (n *Neuron) Run() error {
	var err error
	var cmd string

	n.state = running
	go func() {
		/*
			A dynamic select statement using the reflect package.

			Select executes a select operation described by the list of cases.
			Like the Go select statement, it blocks until at least one of the cases
			can proceed, makes a uniform pseudo-random choice, and then executes that
			case. It returns the index of the chosen case and, if that case was a
			receive operation, the value received and a boolean indicating whether
			the value corresponds to a send on the channel (as opposed to a zero
			value received because the channel is closed).
			You pass in an array of SelectCase structs that identify the channel to
			select on, the direction of the operation, and a value to send in the
			case of a send operation.
		*/

		cases := make([]reflect.SelectCase, len(n.dendriteChans)+1)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(n.cntlChan)}
		for i, ch := range n.dendriteChans {
			cases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}
		chosen, value, _ /*ok*/ := reflect.Select(cases)
		// ok will be true if the channel has not been closed.
		if chosen == 0 {
			cmd = value.String()
			switch cmd {
			case "status":
				log.Printf("Status...%s", n.Status())
			case "stop":
				log.Printf("Stopped neuron: %p", n)
				return
			}
		} else {
			n.sigma += int(value.Int())
		}
		n.lastInChan = n.dendriteChans[chosen]
	}()

	return err
}
