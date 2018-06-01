package contagiongo

import rv "github.com/kentwait/randomvariate"

// TransmissionModel describes the transmission probability and number of
// pathogens that transmits per event. The model may be constant or
// probabilistic.
type TransmissionModel interface {
	// ID returns the ID for this transmission model.
	ModelID() int
	// Name returns the name for this transmission model.
	ModelName() string
	SetModelID(id int)
	SetModelName(name string)
	// TransmissionProb returns the probability that a transmission event
	// occurs between one host and one neighbor (per capita event) occurs.
	TransmissionProb() float64

	// TransmissionSize returns the number of pathogens transmitted given
	// a transmission event occurs.
	TransmissionSize() int
}

type poissonTransmitter struct {
	modelMetadata
	prob float64
	size float64
}

func (s *poissonTransmitter) TransmissionProb() float64 {
	return s.prob
}

func (s *poissonTransmitter) TransmissionSize() int {
	return rv.Poisson(s.size)
}

type constantTransmitter struct {
	modelMetadata
	prob float64
	size int
}

func (s *constantTransmitter) TransmissionProb() float64 {
	return s.prob
}

func (s *constantTransmitter) TransmissionSize() int {
	return s.size
}
