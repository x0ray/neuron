package neuron

import (
	"fmt"
	"log"
	"reflect"
	"time"
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
	started bool
	neurons   []*Neuron
	cntlChans []chan cntlMsg
)

// Neuron - simulated neuron, composed of this dats structure and
//   its own Go routine (when in running state)
type Neuron struct {
	cntlChan      chan cntlMsg     // input from central neuron controller
	dendriteChans []chan int       // dendrite is input from other neurons
	axonChans     []chan int       // axon is output to other neurons
	connected     bool             // connected to other neurons
	state         nstate           // can be: stopped, running, paused
	sigma         int              // sum of dendrite signals
	threshold     int              // level required for axon output
	lastInChan    chan int         // last dendrite chan received
	waitForChan   chan struct{}    // wiat for quiet neuron (no activity over time)
	tickChan      <-chan time.Time // timer chan 
	activity      int              // number of messages received per second
}

// New - creates a new neuron in the stopped state
func New(threshold int) *Neuron {
	n := new(Neuron)
	n.cntlChan = make(chan cntlMsg)
	cntlChans = append(cntlChans, n.cntlChan)
	n.connected = false
	n.state = stopped
	neurons = append(neurons, n)
	n.tickChan = time.After(time.Second)
	n.threshold = threshold	
	n.activity++   // count creation as activity, so WaitForQuiet() will work 
	return n
}

// NewNeurons - create many neurons in a slice
func NewNeurons(num int, threshold int) []*Neuron {
	ns := make([]*Neuron, num)
	for i := 0; i < num; i++ {
		ns[i] = New(threshold)
	}
	return ns
}

// status - print neuron status
//   should only be called from neurons Go routine
func (n *Neuron) status() string {
	var s string
	s = fmt.Sprintf("\nNeuron: %p\n  State..........: %v\n  Dendrite inputs: %d\n  Axon outputs...: %d\n"+
		"  Connected......: %v\n  Sigma..........: %d\n  Threshold......: %d",
		n, n.state, len(n.dendriteChans), len(n.axonChans), n.connected, n.sigma, n.threshold)
	return s
}

// Status - request status of all neurons be displayed
func Status() {
	if started {
		for i, n := range neurons {
			n.cntlChan <- cntlMsg{cmd: "status"}
			log.Printf("Status requested for neuron%d: %p on cntl chan: %p", i, n ,n.cntlChan)
		}
	} else {
		for i, n := range neurons {
			log.Printf("Status for neuron%d: %p", i, n)
			log.Printf("Status...%s", n.status())
		}
	}
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
	if n.state == running {
		if len(n.dendriteChans) > 0 {
			fireChan := n.dendriteChans[0]
			fireChan <- impulse
		}
	}
	return n
}

// FireMultiple - send an impulse to specific neurons
func (n *Neuron) FireMultiple(neurons []*Neuron, impulse int) *Neuron {
	for _, m := range neurons {
		if m.state == running {
			if len(m.dendriteChans) > 0 {
				fireChan := m.dendriteChans[0]
				fireChan <- impulse
			}
		}
	}
	return n
}

// WaitForQuiet - wait until this neuron becomes inactive for a time
func (n *Neuron) WaitForQuiet() {
	if n.waitForChan == nil {
		n.waitForChan = make(chan struct{})
		<-n.waitForChan
	} else {
		<-n.waitForChan
	}
}

// StartAllNeurons - start all the defined neurons
func StartAllNeurons() {	
	for i, n := range neurons {
		n.Run()
		log.Printf("Started neuron%d: %p", i, n)
	} 
	started = true
}

// StopAllNeurons - start all the defined neurons
func StopAllNeurons() {
	for i, n := range cntlChans {
		n <- cntlMsg{cmd: "stop"}
		log.Printf("Stopped neuron%d on cntl chan: %p", i, n)
	}
}

// Run - start all neuron channels executing
func (n *Neuron) Run() error {
	var err error
	var cmd string

	if n.state == stopped {
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
			const fixedCases = 2
			cases := make([]reflect.SelectCase, len(n.dendriteChans)+fixedCases)
			cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(n.cntlChan)}
			cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(n.tickChan)}
			for i, ch := range n.dendriteChans {
				cases[i+fixedCases] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
			}
			chosen, value, _ /*ok*/ := reflect.Select(cases)
			// ok will be true if the channel has not been closed.
			if chosen == 0 { // cntlChan selected
				cmd = value.String()
				switch cmd {
				case "status":
					log.Printf("Status...%s", n.status())
				case "stop":
					log.Printf("Stopped neuron: %p", n)
					return
				}
			} else if chosen == 1 { // timer chan selected
				if n.activity > 0 {
					if n.waitForChan != nil {
						n.waitForChan <- struct{}{}
					}
					n.activity = 0
				}
			} else { // dendrite chan selected
				n.activity++
				n.sigma += int(value.Int())
				if n.sigma > n.threshold { // trigger the axom ?
					for _, m := range n.axonChans {
						m <- 1
					}
				}
			}
			if len(n.dendriteChans) > 0 {
				n.lastInChan = n.dendriteChans[chosen]
			}
		}()
		n.state = running
	}
	return err
}
