package produce

type IProduce interface {
	SendMessage() error
}
