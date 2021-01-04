package produce

type Kafka struct {
	Span *Span
}

func (span Kafka) SendMessage1() error {
	return nil
}
