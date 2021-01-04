package produce

type Redis struct {
	Span *Span
}

func (span Redis) SendMessage() error {
	return nil
}
