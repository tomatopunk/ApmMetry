package produce

type KSql struct {
	Span *Span
}

func (span KSql) SendMessage() error {
	return nil
}
