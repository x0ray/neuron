# Package Neuron

Simulated concurrent neurons in Go

## Introduction  

The neuron package creates simulated neurons. Each neuron exists as a Go routine. The neurons inputs known as dendrites in real neurons are simulated by a Go language channel. The neurons outputs or the axon are also simulated by channels. Additionally each neuron has a control channel which is used for global control of all the neurons, for example to shut al the neurons down.

A critical and less well known Go language feature that this package requires is the use of the reflect package to create a dynamic select statement which allows the neural network to be connected and modified at run time.

## Methods

#### New - creates a new neuron in the stopped state
``` go
func New(threshold int) *Neuron 
```

#### NewNeurons - create many neurons in a slice
``` go
func NewNeurons(num int, threshold int) []*Neuron 
```

#### Status - print neuron status
``` go
func (n *Neuron) Status() string 
```

#### Link - connect an axom channel of a neuron to a dendrite channel of a neuron
``` go
func (n *Neuron) Link(to *Neuron) *Neuron
```

#### LinkManyToOne - connect an axom channel of many neurons to dendrite channels of a neuron
``` go
func (n *Neuron) LinkManyToOne(many []*Neuron) *Neuron 
```

#### LinkOneToMany - connect an axom channel of one neuron to dendrite channels of many neurons
``` go
func (n *Neuron) LinkOneToMany(many []*Neuron) *Neuron 
```

#### Fire - send an impulse to a specific neuron
``` go
func (n *Neuron) Fire(impulse int) *Neuron 
```

#### FireMultiple - send an impulse to specific neurons
``` go
func (n *Neuron) FireMultiple(neurons []*Neuron, impulse int) *Neuron
```

#### StartAllNeurons - start all the defined neurons
``` go
func StartAllNeurons() 
```

#### StopAllNeurons - start all the defined neurons
``` go
func StopAllNeurons() 
```

## Example

This example shows how to establish neurons, and connect them in a network. Then how to start the neural network, print its status and then shut it down. Later examples will demonstrate how to fire the neurons.

### Go test program
``` go
	// make neurons
	in := u.New()
	plain := u.NewNeurons(10)
	out := u.New()

	// connect neurons
	in.LinkOneToMany(plain)
	out.LinkManyToOne(plain)

	// start all neurons
	u.StartAllNeurons()

	// print status of all neurons
	l.Printf(in.Status())
	for _, n := range plain {
		l.Printf(n.Status())
	}
	l.Printf(out.Status())

	// stop all neurons
	u.StopAllNeurons()
```

### Output
``` sh
$ neurotest
2021/08/26 11:49:47 neurotest version: 0.0.1 Started
2021/08/26 11:49:47 Started neuron0: 0xc0000281e0
2021/08/26 11:49:47 Started neuron1: 0xc000028240
2021/08/26 11:49:47 Started neuron2: 0xc0000282a0
2021/08/26 11:49:47 Started neuron3: 0xc000028300
2021/08/26 11:49:47 Started neuron4: 0xc000028360
2021/08/26 11:49:47 Started neuron5: 0xc0000283c0
2021/08/26 11:49:47 Started neuron6: 0xc000028420
2021/08/26 11:49:47 Started neuron7: 0xc000028480
2021/08/26 11:49:47 Started neuron8: 0xc0000284e0
2021/08/26 11:49:47 Started neuron9: 0xc000028540
2021/08/26 11:49:47 Started neuron10: 0xc0000285a0
2021/08/26 11:49:47 Started neuron11: 0xc000028600
2021/08/26 11:49:47 
Neuron: 0xc0000281e0
  State..........: running
  Dendrite inputs: 0
  Axon outputs...: 10
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028240
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc0000282a0
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028300
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028360
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc0000283c0
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028420
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028480
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc0000284e0
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028540
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc0000285a0
  State..........: running
  Dendrite inputs: 1
  Axon outputs...: 1
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 
Neuron: 0xc000028600
  State..........: running
  Dendrite inputs: 10
  Axon outputs...: 0
  Connected......: true
  Sigma..........: 0
  Threshold......: 0
2021/08/26 11:49:47 Stopped neuron0 on cntl chan: 0xc0000221e0
2021/08/26 11:49:47 Stopped neuron1 on cntl chan: 0xc000022240
2021/08/26 11:49:47 Stopped neuron2 on cntl chan: 0xc0000222a0
2021/08/26 11:49:47 Stopped neuron3 on cntl chan: 0xc000022300
2021/08/26 11:49:47 Stopped neuron4 on cntl chan: 0xc000022360
2021/08/26 11:49:47 Stopped neuron5 on cntl chan: 0xc0000223c0
2021/08/26 11:49:47 Stopped neuron6 on cntl chan: 0xc000022420
2021/08/26 11:49:47 Stopped neuron7 on cntl chan: 0xc000022480
2021/08/26 11:49:47 Stopped neuron8 on cntl chan: 0xc0000224e0
2021/08/26 11:49:47 Stopped neuron9 on cntl chan: 0xc000022540
2021/08/26 11:49:47 Stopped neuron10 on cntl chan: 0xc0000225a0
2021/08/26 11:49:47 Stopped neuron11 on cntl chan: 0xc000022600
2021/08/26 11:49:47 neurotest Ended
```