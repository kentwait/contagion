package contagiongo

type DataLogger interface {
	WriteGenotypes()
	WriteGenotypeFreq()
	WriteStatus()
	WriteGenotypeHistory()
}
